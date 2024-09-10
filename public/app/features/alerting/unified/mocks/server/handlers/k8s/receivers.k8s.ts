import { camelCase } from 'lodash';
import { HttpResponse, http } from 'msw';

import alertmanagerConfig from 'app/features/alerting/unified/components/contact-points/__mocks__/alertmanager.config.mock.json';
import { ALERTING_API_SERVER_BASE_URL, getK8sResponse } from 'app/features/alerting/unified/mocks/server/utils';
import { ComGithubGrafanaGrafanaPkgApisAlertingNotificationsV0Alpha1Receiver } from 'app/features/alerting/unified/openapi/receiversApi.gen';
import { PROVENANCE_ANNOTATION, PROVENANCE_NONE } from 'app/features/alerting/unified/utils/k8s/constants';
import { ANNOTATION_PREFIX_ACCESS } from 'app/features/alerting/unified/utils/k8s/utils';
import { AlertManagerCortexConfig } from 'app/plugins/datasource/alertmanager/types';

const config: AlertManagerCortexConfig = alertmanagerConfig;

// Turn our mock alertmanager config into the format that we expect to be returned by the k8s API
const mappedReceivers =
  config.alertmanager_config?.receivers?.map((contactPoint) => {
    const provenance =
      contactPoint.grafana_managed_receiver_configs?.find((integration) => {
        return integration.provenance;
      })?.provenance || PROVENANCE_NONE;
    return {
      metadata: {
        // This isn't exactly accurate, but its the cleanest way to use the same data for AM config and K8S responses
        uid: camelCase(contactPoint.name),
        annotations: {
          [PROVENANCE_ANNOTATION]: provenance,
          [`${ANNOTATION_PREFIX_ACCESS}/canAdmin`]: 'true',
          [`${ANNOTATION_PREFIX_ACCESS}/canWrite`]: 'true',
          [`${ANNOTATION_PREFIX_ACCESS}/canDelete`]: 'true',
        },
      },
      spec: {
        title: contactPoint.name,
        integrations: contactPoint.grafana_managed_receiver_configs || [],
      },
    };
  }) || [];

const parsedReceivers = getK8sResponse<ComGithubGrafanaGrafanaPkgApisAlertingNotificationsV0Alpha1Receiver>(
  'ReceiverList',
  mappedReceivers
);

const listNamespacedReceiverHandler = () =>
  http.get<{ namespace: string }>(`${ALERTING_API_SERVER_BASE_URL}/namespaces/:namespace/receivers`, () => {
    return HttpResponse.json(parsedReceivers);
  });

const createNamespacedReceiverHandler = () =>
  http.post<{ namespace: string }>(
    `${ALERTING_API_SERVER_BASE_URL}/namespaces/:namespace/receivers`,
    async ({ request }) => {
      const body = await request.json();
      return HttpResponse.json(body);
    }
  );

const deleteNamespacedReceiverHandler = () =>
  http.delete<{ namespace: string; name: string }>(
    `${ALERTING_API_SERVER_BASE_URL}/namespaces/:namespace/receivers/:name`,
    ({ params }) => {
      const { name } = params;
      const matchedReceiver = parsedReceivers.items.find((receiver) => receiver.metadata.uid === name);
      if (matchedReceiver) {
        return HttpResponse.json(parsedReceivers);
      }
      return HttpResponse.json({}, { status: 404 });
    }
  );

const handlers = [
  listNamespacedReceiverHandler(),
  createNamespacedReceiverHandler(),
  deleteNamespacedReceiverHandler(),
];
export default handlers;

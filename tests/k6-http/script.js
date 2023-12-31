/*
 * Pedágio
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * OpenAPI spec version: 1.0
 *
 * NOTE: This class is auto generated by OpenAPI Generator.
 * https://github.com/OpenAPITools/openapi-generator
 *
 * OpenAPI generator version: 7.1.0-SNAPSHOT
 */


import http from "k6/http";
import { group, check, sleep } from "k6";

// @TODO ajustar o esquema da URL em função da disponibilidade de TLS
const BASE_URL = `http://${__ENV.HOSTNAME}`;

// Sleep duration between successive requests.
// You might want to edit the value of this variable or remove calls to the sleep function on the script.
const SLEEP_DURATION = 0.1;

// Global variables should be initialized.

export default function() {
    // @TODO ajustar a URL do service em função da especificação
    group("/", () => {
        // Request No. 1: TollStationService_OpenGate
        {
            let url = BASE_URL + `/`;

            // TODO: edit the parameters of the request body.
            let body = {"tag": "tag-1"};
            let params = {headers: {"Content-Type": "application/json", "Accept": "application/json"}};
            let request = http.post(url, JSON.stringify(body), params);

            check(request, {
                "A successful response.": (r) => r.status === 200
            });
        }
    });

}

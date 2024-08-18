/* tslint:disable */
/* eslint-disable */
/**
 * Vibrain API
 * This is a simple API for Vibrain project.
 *
 * The version of the OpenAPI document: 1.0
 * Contact: vibrain@vaayne.com
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { mapValues } from "../runtime";
/**
 *
 * @export
 * @interface HttpserverCreateThreadMessageRequest
 */
export interface HttpserverCreateThreadMessageRequest {
  /**
   *
   * @type {string}
   * @memberof HttpserverCreateThreadMessageRequest
   */
  model?: string;
  /**
   *
   * @type {string}
   * @memberof HttpserverCreateThreadMessageRequest
   */
  role?: string;
  /**
   *
   * @type {string}
   * @memberof HttpserverCreateThreadMessageRequest
   */
  text?: string;
}

/**
 * Check if a given object implements the HttpserverCreateThreadMessageRequest interface.
 */
export function instanceOfHttpserverCreateThreadMessageRequest(
  value: object,
): value is HttpserverCreateThreadMessageRequest {
  return true;
}

export function HttpserverCreateThreadMessageRequestFromJSON(
  json: any,
): HttpserverCreateThreadMessageRequest {
  return HttpserverCreateThreadMessageRequestFromJSONTyped(json, false);
}

export function HttpserverCreateThreadMessageRequestFromJSONTyped(
  json: any,
  ignoreDiscriminator: boolean,
): HttpserverCreateThreadMessageRequest {
  if (json == null) {
    return json;
  }
  return {
    model: json["model"] == null ? undefined : json["model"],
    role: json["role"] == null ? undefined : json["role"],
    text: json["text"] == null ? undefined : json["text"],
  };
}

export function HttpserverCreateThreadMessageRequestToJSON(
  value?: HttpserverCreateThreadMessageRequest | null,
): any {
  if (value == null) {
    return value;
  }
  return {
    model: value["model"],
    role: value["role"],
    text: value["text"],
  };
}
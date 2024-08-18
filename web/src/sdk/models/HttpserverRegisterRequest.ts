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
 * @interface HttpserverRegisterRequest
 */
export interface HttpserverRegisterRequest {
  /**
   *
   * @type {string}
   * @memberof HttpserverRegisterRequest
   */
  email?: string;
  /**
   *
   * @type {string}
   * @memberof HttpserverRegisterRequest
   */
  password?: string;
  /**
   *
   * @type {string}
   * @memberof HttpserverRegisterRequest
   */
  username?: string;
}

/**
 * Check if a given object implements the HttpserverRegisterRequest interface.
 */
export function instanceOfHttpserverRegisterRequest(
  value: object,
): value is HttpserverRegisterRequest {
  return true;
}

export function HttpserverRegisterRequestFromJSON(
  json: any,
): HttpserverRegisterRequest {
  return HttpserverRegisterRequestFromJSONTyped(json, false);
}

export function HttpserverRegisterRequestFromJSONTyped(
  json: any,
  ignoreDiscriminator: boolean,
): HttpserverRegisterRequest {
  if (json == null) {
    return json;
  }
  return {
    email: json["email"] == null ? undefined : json["email"],
    password: json["password"] == null ? undefined : json["password"],
    username: json["username"] == null ? undefined : json["username"],
  };
}

export function HttpserverRegisterRequestToJSON(
  value?: HttpserverRegisterRequest | null,
): any {
  if (value == null) {
    return value;
  }
  return {
    email: value["email"],
    password: value["password"],
    username: value["username"],
  };
}

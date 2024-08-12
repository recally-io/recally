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
 * JSONResult represents the structure of the JSON response.
 * @export
 * @interface HandlersJSONResult
 */
export interface HandlersJSONResult {
  /**
   * Code is an integer value that represents the HTTP status code.
   * @type {number}
   * @memberof HandlersJSONResult
   */
  code?: number;
  /**
   * Data is an interface value that represents the data of the response.
   * @type {object}
   * @memberof HandlersJSONResult
   */
  data?: object;
  /**
   * Error is an error value that represents the error of the response.
   * @type {object}
   * @memberof HandlersJSONResult
   */
  error?: object;
  /**
   * Message is a string value that represents the message of the response.
   * @type {string}
   * @memberof HandlersJSONResult
   */
  message?: string;
  /**
   * Success is a boolean value that indicates whether the request was successful.
   * @type {boolean}
   * @memberof HandlersJSONResult
   */
  success?: boolean;
}

/**
 * Check if a given object implements the HandlersJSONResult interface.
 */
export function instanceOfHandlersJSONResult(
  value: object,
): value is HandlersJSONResult {
  return true;
}

export function HandlersJSONResultFromJSON(json: any): HandlersJSONResult {
  return HandlersJSONResultFromJSONTyped(json, false);
}

export function HandlersJSONResultFromJSONTyped(
  json: any,
  ignoreDiscriminator: boolean,
): HandlersJSONResult {
  if (json == null) {
    return json;
  }
  return {
    code: json["code"] == null ? undefined : json["code"],
    data: json["data"] == null ? undefined : json["data"],
    error: json["error"] == null ? undefined : json["error"],
    message: json["message"] == null ? undefined : json["message"],
    success: json["success"] == null ? undefined : json["success"],
  };
}

export function HandlersJSONResultToJSON(
  value?: HandlersJSONResult | null,
): any {
  if (value == null) {
    return value;
  }
  return {
    code: value["code"],
    data: value["data"],
    error: value["error"],
    message: value["message"],
    success: value["success"],
  };
}

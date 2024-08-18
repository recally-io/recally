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
 * @interface ToolsWebSummaryGet200Response
 */
export interface ToolsWebSummaryGet200Response {
  /**
   * Code is an integer value that represents the HTTP status code.
   * @type {number}
   * @memberof ToolsWebSummaryGet200Response
   */
  code?: number;
  /**
   *
   * @type {string}
   * @memberof ToolsWebSummaryGet200Response
   */
  data?: string;
  /**
   * Error is an error value that represents the error of the response.
   * @type {object}
   * @memberof ToolsWebSummaryGet200Response
   */
  error?: object;
  /**
   * Message is a string value that represents the message of the response.
   * @type {string}
   * @memberof ToolsWebSummaryGet200Response
   */
  message?: string;
  /**
   * Success is a boolean value that indicates whether the request was successful.
   * @type {boolean}
   * @memberof ToolsWebSummaryGet200Response
   */
  success?: boolean;
}

/**
 * Check if a given object implements the ToolsWebSummaryGet200Response interface.
 */
export function instanceOfToolsWebSummaryGet200Response(
  value: object,
): value is ToolsWebSummaryGet200Response {
  return true;
}

export function ToolsWebSummaryGet200ResponseFromJSON(
  json: any,
): ToolsWebSummaryGet200Response {
  return ToolsWebSummaryGet200ResponseFromJSONTyped(json, false);
}

export function ToolsWebSummaryGet200ResponseFromJSONTyped(
  json: any,
  ignoreDiscriminator: boolean,
): ToolsWebSummaryGet200Response {
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

export function ToolsWebSummaryGet200ResponseToJSON(
  value?: ToolsWebSummaryGet200Response | null,
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

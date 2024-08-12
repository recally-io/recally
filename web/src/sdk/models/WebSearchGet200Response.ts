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
import type { JinasearcherContent } from "./JinasearcherContent";
import {
  JinasearcherContentFromJSON,
  JinasearcherContentFromJSONTyped,
  JinasearcherContentToJSON,
} from "./JinasearcherContent";

/**
 *
 * @export
 * @interface WebSearchGet200Response
 */
export interface WebSearchGet200Response {
  /**
   * Code is an integer value that represents the HTTP status code.
   * @type {number}
   * @memberof WebSearchGet200Response
   */
  code?: number;
  /**
   *
   * @type {JinasearcherContent}
   * @memberof WebSearchGet200Response
   */
  data?: JinasearcherContent;
  /**
   * Error is an error value that represents the error of the response.
   * @type {object}
   * @memberof WebSearchGet200Response
   */
  error?: object;
  /**
   * Message is a string value that represents the message of the response.
   * @type {string}
   * @memberof WebSearchGet200Response
   */
  message?: string;
  /**
   * Success is a boolean value that indicates whether the request was successful.
   * @type {boolean}
   * @memberof WebSearchGet200Response
   */
  success?: boolean;
}

/**
 * Check if a given object implements the WebSearchGet200Response interface.
 */
export function instanceOfWebSearchGet200Response(
  value: object,
): value is WebSearchGet200Response {
  return true;
}

export function WebSearchGet200ResponseFromJSON(
  json: any,
): WebSearchGet200Response {
  return WebSearchGet200ResponseFromJSONTyped(json, false);
}

export function WebSearchGet200ResponseFromJSONTyped(
  json: any,
  ignoreDiscriminator: boolean,
): WebSearchGet200Response {
  if (json == null) {
    return json;
  }
  return {
    code: json["code"] == null ? undefined : json["code"],
    data:
      json["data"] == null
        ? undefined
        : JinasearcherContentFromJSON(json["data"]),
    error: json["error"] == null ? undefined : json["error"],
    message: json["message"] == null ? undefined : json["message"],
    success: json["success"] == null ? undefined : json["success"],
  };
}

export function WebSearchGet200ResponseToJSON(
  value?: WebSearchGet200Response | null,
): any {
  if (value == null) {
    return value;
  }
  return {
    code: value["code"],
    data: JinasearcherContentToJSON(value["data"]),
    error: value["error"],
    message: value["message"],
    success: value["success"],
  };
}
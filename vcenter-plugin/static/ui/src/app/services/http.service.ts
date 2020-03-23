import {Injectable} from "@angular/core";
import {Headers, Http, Response} from "@angular/http";
import {Observable} from "rxjs/Observable";
import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/share';
import 'rxjs/add/operator/mergeMap';
import 'rxjs/add/observable/throw';
import {GlobalService, SessionInfo} from "./global.service";

import * as _ from 'lodash';
import { RequestOptionsArgs } from "../../../node_modules/@angular/http/src/interfaces";

@Injectable()
export class HttpService extends GlobalService {
   private readonly JSON_HEADERS = {
      'Content-Type': 'application/json;charset=utf-8',
      'Accept': 'application/json;charset=utf-8'
   };

   private readonly CACHE_CONTROL_HEADERS = {
      'Cache-Control': 'no-cache',
      'Pragma': 'no-cache',
      'Expires': 'Sat, 01 Jan 2000 00:00:00 GMT'
   };

   constructor(private http: Http) {
      super();
   }

   private getSessionInfo(): Observable<SessionInfo> {
      return new Observable(observer => {
         this.htmlClientSdk.app.getSessionInfo((sessionInfo: SessionInfo) => {
            observer.next(sessionInfo);
            observer.complete();
         })
      });
   }

   /**
    * Performs HTTP GET request against the remote plugin REST endpoint.
    *
    * @param {string} path
    *    Sub-path the GET request will be made for.
    */
   get(path: string, params?: any): Observable<any> {
      return this.getSessionInfo().flatMap((si: SessionInfo) => {
         return this.doGetOrDelete(
               this.http.get.bind(this.http), si, path, params);
      });
   }

   /**
    * Performs HTTP POST request against the remote plugin REST endpoint.
    *
    * @param {string} path
    *    Sub-path the request will be made for.
    * @param {any} body (optional)
    *    The body of the request.
    * @param {string} params optional
    *    Query string parameters of the request.
    */
   post(path: string, body?: any, params?: string): Observable<any> {
      return this.getSessionInfo().flatMap((si: SessionInfo) => {
         return this.doPostOrPut(
               this.http.post.bind(this.http), si, path, body, params);
      });
   }

   /**
    * Performs HTTP PUT request against the remote plugin REST endpoint.
    *
    * @param {string} path
    *    Sub-path the request will be made for.
    * @param {any} body (optional)
    *    The body of the request.
    * @param {string} params optional
    *    Query string parameters of the request.
    */
   put(path: string, body?: any, params?: string): Observable<any> {
      return this.getSessionInfo().flatMap((si: SessionInfo) => {
         return this.doPostOrPut(
               this.http.put.bind(this.http), si, path, body, params);
      });
   }

   /**
    * Performs HTTP DELETE request against the remote plugin REST endpoint.
    *
    * @param {string} path
    *    Sub-path the request will be made for.
    * @param {string} params optional
    *    Query string parameters of the request.
    */
   delete(path: string, params?: string): Observable<any> {
      return this.getSessionInfo().flatMap((si: SessionInfo) => {
         return this.doGetOrDelete(
               this.http.delete.bind(this.http), si, path, params);
      });
   }

   protected static convertPath(path: string): string {
      return `${GlobalService.WEB_CONTEXT_PATH}/${path}`;
   }

   private doGetOrDelete(httpMethod: Function, si: SessionInfo, path: string,
                         params?: string): Observable<any> {
      return this.getSessionInfo().flatMap(() => {
         return httpMethod(
               HttpService.convertPath(path),
               this.getRequestOptionsArgs(params, si))
               .catch(this.handleError)
               .share()
               .map((response: Response) => HttpService.parseResponse(response));
      });
   }

   private doPostOrPut(httpMethod: Function, si: SessionInfo, path: string,
                       body?: any, params?: string): Observable<any> {
      return this.getSessionInfo().flatMap(() => {
         return httpMethod(
               HttpService.convertPath(path),
               body,
               this.getRequestOptionsArgs(params, si))
               .catch(this.handleError)
               .share()
               .map((response: Response) => HttpService.parseResponse(response));
      });
   }

   private getRequestOptionsArgs(params: string,
                                 si: SessionInfo): RequestOptionsArgs {
      return <RequestOptionsArgs>
            {
               params: params,
               headers: new Headers(_.extend(
                     {},
                     this.JSON_HEADERS,
                     this.CACHE_CONTROL_HEADERS,
                     {
                        "vmware-api-gateway-url": this.htmlClientSdk.app
                              .getApiEndpoints().uiApiEndpoint.fullUrl,
                        "vmware-api-session-id": si.sessionToken
                     }
               ))
            };
   }

   /**
    * Default callback for handling errors on http calls.
    */
   private handleError = (error: any) => {
      const DEFAULT_ERROR_MESSAGE = 'Backend server error';
      let throwable: any;

      try {
         throwable = (error instanceof Response) ? error.json() : error;
      } catch (e) {
         // Response is not JSON, so expect a general HTTP error.
         throwable = error;
      }
      return Observable.throw(throwable || DEFAULT_ERROR_MESSAGE);
   };

   /**
    * Parses HTTP response as JSON.
    *
    * @param {Response} response
    *    HTTP response.
    *
    * @throw Error
    *    If the response is not valid JSON.
    */
   private static parseResponse(response: Response): any {
      let out: any = null;
      try {
         out = response.json();
      } catch (e) {
         Observable.throw(e);
      }
      return out;
   }
}
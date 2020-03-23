/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Injectable} from "@angular/core";

import "rxjs/add/operator/toPromise";
import {Observable} from "rxjs/Observable";

export interface SessionInfo {
   sessionToken: string;
}

@Injectable()
export class GlobalService {
   protected static WEB_CONTEXT_PATH = `${window.location.origin}${window.location.pathname}/../rest`;
   public readonly htmlClientSdk: any;

   constructor() {
      this.htmlClientSdk = (<any>window).htmlClientSdk;
   }
}

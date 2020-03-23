/* Copyright (c) 2017 VMware, Inc. All rights reserved. */
import {Injectable} from '@angular/core';

/**
 * In production mode there is no 'window' object in build-time.
 * This wrapper-service solves the problem
 */
@Injectable()
export class WindowService {
   get window(): any {
      return (() => {
         return window;
      })();
   }

}
/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Injectable} from "@angular/core";

@Injectable()
export class ValidationUtil {
   public static isNullOrEmpty(name: string) {
      return !(name && name.trim().length != 0);
   }
}

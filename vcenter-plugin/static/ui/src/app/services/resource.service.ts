/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Injectable} from "@angular/core";
import {TranslateService} from "@ngx-translate/core";

@Injectable()
export class ResourceService {

   private localizedStrings: Object;

   constructor(private translate: TranslateService) {
      let strings = ["shared.modal.createChassis",
         "shared.modal.editChassis",
         "common.active",
         "common.standBy",
         "shared.modal.createChassis",
         "shared.modal.editChassis",
         "actions.create.emptyNameError",
         "actions.create.usedNameError"];
      this.translate.get(strings).subscribe(
         (result: Object) => {
            this.localizedStrings = result;
         });
   }

   public getString(str: string): string {
      return this.localizedStrings && this.localizedStrings.hasOwnProperty(str) ? this.localizedStrings[str] : str;
   }
}

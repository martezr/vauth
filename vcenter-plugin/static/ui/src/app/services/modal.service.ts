/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Injectable} from "@angular/core";
import {ResourceService} from "./resource.service";

export interface ModalConfig {
   url: string;
   title?: string;
   size?: {
      width: number,
      height: number
   };
   closable?: boolean;
   onClosed?: (result?: any) => void;
   contextObjects?: any[];
   customData?: any;
}

@Injectable()
export class ModalConfigService {

   constructor(private resources: ResourceService){
   }

   public createAddConfig() {
      let addAction: ModalConfig = {
         url: "index.html?view=create",
         title: this.resources.getString("shared.modal.createChassis"),
         size: {width: 700, height: 400}
      };
      return addAction;
   }

   public createAddWizardConfig() {
      let addWizardAction: ModalConfig = {
         url: "index.html?view=create-wizard",
         closable: false,
         size: {width: 800, height: 500}
      };
      return addWizardAction;
   }

   public createDeleteConfig() {
      let deleteAction: ModalConfig = {
         url: "index.html?view=delete",
         size: {width: 400, height: 120},
         closable: false
      };
      return deleteAction;
   }

   public createEditConfig() {
      let editAction: ModalConfig = {
         url: "index.html?view=edit",
         title:  this.resources.getString("shared.modal.editChassis"),
         size: {width: 780, height: 500}
      };
      return editAction;
   }
}

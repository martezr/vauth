/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Inject, OnInit} from '@angular/core';
import {Chassis} from "../../../model/chassis.model";
import {ChassisService} from "../../../services/chassis.service";
import {ResourceService} from "../../../services/resource.service";
import {ActivatedRoute} from '@angular/router';
import {ValidationUtil} from "../../../services/validation";

export abstract class CreateEditChassis implements OnInit{
   public chassis: Chassis;
   public chassisExists: boolean;
   public action: string;

   static readonly EDIT_ACTION: string = "edit";

   constructor(@Inject(ChassisService) public chassisService: ChassisService,
               @Inject(ResourceService) public resourceService: ResourceService,
               @Inject(ActivatedRoute) private route: ActivatedRoute) {
      this.action = route.snapshot.url[0].path;
      this.chassis = new Chassis();
   }

   ngOnInit(): void {
      if (this.isEditAction()) {
         (<any>Object).assign(this.chassis,
            this.chassisService.htmlClientSdk.app.getContextObjects()[0]);
      }
   }

   abstract onSubmit(): void;

   onCancel(result?: any): void {
      this.chassisService.htmlClientSdk.modal.close(result);
   }

   isValid(): boolean{
      return !ValidationUtil.isNullOrEmpty(this.chassis.name);
   }

   invalidNameError(): string {
      if (!this.chassisExists) {
         return this.resourceService.getString("actions.create.emptyNameError");
      } else {
         return this.resourceService.getString("actions.create.usedNameError");
      }
   }

   isEditAction(): boolean {
      return this.action === CreateEditChassis.EDIT_ACTION;
   }

   nameChanged(newValue: string): void {
      this.chassis.name = newValue;
   }

   serverTypeChanged(newValue: string): void {
      this.chassis.serverType = newValue;
   }

   dimensionsChanged(newValue: string): void {
      this.chassis.dimensions = newValue;
   }

   onCreateChassisFailed(): void {
      this.chassisExists = true;
   }

   create(): void {
      this.chassisExists = false;

      this.chassisService.create(this.chassis)
         .then((result: Chassis) => {
            if (result) {
               this.onCancel(result);
            } else {
               this.onCreateChassisFailed();
            }
         })
         .catch(() => {
            this.onCancel();
         });
   }

   edit(): void {
      this.chassisService.edit(this.chassis)
         .then((chassisObjectUpdated: boolean) => {
            if (!chassisObjectUpdated) {
               this.chassisExists = true;
            } else {
               this.onCancel(this.chassis);
            }
         })
         .catch(() => {
            this.onCancel();
         });
   }

   onHostsSelectionChange(selectedHosts: string[]): void {
      this.chassis.relatedHostsIds = selectedHosts;
   }

   onNavigateToHostObject(): void {
      this.onCancel();
   }
}

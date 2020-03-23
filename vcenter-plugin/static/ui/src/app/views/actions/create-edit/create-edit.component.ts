/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import { Component } from '@angular/core';
import { CreateEditChassis } from "./create-edit";

/**
 * Represents a form for creating or editing a chassis.
 */
@Component(
      {
         templateUrl: './create-edit.component.html'
      }
)
export class CreateEditComponent extends CreateEditChassis {

   public onSubmit(): void {
      if (this.isEditAction()) {
         this.edit();
      } else {
         this.create();
      }
   }
}

/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {
   Component, Input,
   Output, EventEmitter
} from '@angular/core';
import {Chassis} from "../../model/chassis.model";

@Component({
   selector: "details-view",
   templateUrl: './details-view.component.html',
   styleUrls: ["./details-view.component.scss"]
})

export class DetailsViewComponent {

   @Input()
   chassis: Chassis;

   @Output()
   chassisUpdated = new EventEmitter<Chassis>();

   @Output()
   chassisDeleted = new EventEmitter<Chassis>();

   onChassisUpdated(chassis: Chassis): void {
      this.chassis = chassis;
      this.chassisUpdated.emit(chassis);
   }

   onChassisDeleted(chassis: Chassis): void {
      this.chassisDeleted.emit(chassis);
   }
}

/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {
   Component,
   Input,
   OnInit,
   NgZone,
   Output,
   EventEmitter
} from '@angular/core';

import {Chassis} from "../../model/chassis.model";
import {ChassisService} from "../../services/chassis.service";
import {ModalConfig, ModalConfigService} from "../../services/modal.service";

@Component({
   selector: "custom-header",
   templateUrl: './header.component.html',
   styleUrls: ["./header.component.scss"],
})

export class HeaderComponent implements OnInit {

   @Input()
   chassis: Chassis;

   @Output()
   chassisUpdated = new EventEmitter<Chassis>();

   @Output()
   chassisDeleted = new EventEmitter<Chassis>();

   showAlert: boolean;

   constructor(private zone: NgZone, private chassisService: ChassisService,
         private modalService: ModalConfigService) {
   }

   ngOnInit(): void {
      this.showAlert = false;
   }

   onDelete(): void {
      let config: ModalConfig = this.modalService.createDeleteConfig();
      config.contextObjects = [this.chassis.id];
      config.onClosed = (result: boolean) => {
         if (result) {
            this.zone.run(() => {
               this.chassisDeleted.emit(this.chassis);
            });
         }
      };
      this.chassisService.htmlClientSdk.modal.open(config);
   }

   onEdit(): void {
      let config: ModalConfig = this.modalService.createEditConfig();
      config.contextObjects = [Object.assign(new Chassis(), this.chassis)];
      config.onClosed = (result: Chassis | null | undefined) => {
         if (result) {
            this.zone.run(() => {
               this.chassis = result;
               this.chassisUpdated.emit(this.chassis);
            });
         }
      };
      this.chassisService.htmlClientSdk.modal.open(config);
   }

   onActivate(): void {
      let newChassis = Object.assign(new Chassis(), this.chassis);
      newChassis.isActive = true;
      this.chassisService.edit(newChassis)
            .then(() => {
               this.showAlert = true;

               this.chassis = Object.assign(new Chassis(), this.chassis);
               this.chassis.isActive = true;

               this.chassisUpdated.emit(this.chassis);
            })
            .catch(() => {
            });
   }

   onAlertClose(): void {
      this.showAlert = false;
   }
}

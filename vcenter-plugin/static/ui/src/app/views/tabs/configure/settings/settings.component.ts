/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Component, Input, Output, NgZone, EventEmitter} from '@angular/core';
import {ChassisService} from "../../../../services/chassis.service";
import {ModalConfig, ModalConfigService} from "../../../../services/modal.service";
import {Chassis} from "../../../../model/chassis.model";

@Component(
   {
      selector: 'settings-view',
      templateUrl: './settings.component.html'
   }
)

export class SettingsComponent {

   constructor(private chassisService: ChassisService,
               private zone: NgZone,
               private modalService: ModalConfigService) {
   }

   @Input()
   chassis: Chassis;


   @Output()
   chassisUpdated = new EventEmitter<Chassis>();

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
}

/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Component, Input, Output, EventEmitter} from '@angular/core';
import {Chassis} from "../../../model/chassis.model";

@Component(
   {
      selector: 'monitor-view',
      templateUrl: './monitor.component.html'
   }
)

export class MonitorComponent{
   @Input()
   chassis: Chassis;

   @Output()
   chassisUpdated = new EventEmitter<Chassis>();

    refreshData(): void {
        this.chassis = Object.assign(new Chassis(), this.chassis);
        this.chassis.healthStatus = Math.floor(Math.random() * 100) + 1;
        this.chassis.complianceStatus = Math.floor(Math.random() * 100) + 1;

        this.chassisUpdated.emit(this.chassis);
    }
}

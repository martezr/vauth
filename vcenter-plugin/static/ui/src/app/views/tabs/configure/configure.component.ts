/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Component, Input, Output, EventEmitter, OnInit} from '@angular/core';
import {Chassis} from "../../../model/chassis.model";

@Component(
   {
      selector: 'configure-view',
      templateUrl: './configure.component.html',
      styleUrls: ['./configure.component.scss'],
   }
)

export class ConfigureComponent implements OnInit {
   @Input()
   chassis: Chassis;

   @Output()
   chassisUpdated = new EventEmitter<Chassis>();

   contentType: string;

   ngOnInit():void {
       this.contentType = '';
   }

   setContent(name: string): void {
      this.contentType = name;
   }

   afterChassisUpdatedHandler(chassis: Chassis): void {
      this.chassis = chassis;
      this.chassisUpdated.emit(chassis);
   }
}

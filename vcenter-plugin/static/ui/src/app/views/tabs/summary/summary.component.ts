/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Component, Input} from '@angular/core';
import {Chassis} from "../../../model/chassis.model";

@Component(
   {
      selector: 'summary-view',
      templateUrl: './summary.component.html'
   }
)

export class SummaryComponent {
    @Input()
    chassis: Chassis;
}

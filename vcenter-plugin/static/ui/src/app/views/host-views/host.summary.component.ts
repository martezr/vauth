/* Copyright (c) 2018 VMware, Inc. All rights reserved. */
import { Component, Inject, OnInit } from "@angular/core";
import { ChassisService } from "../../services/chassis.service";
import { Chassis } from "../../model/chassis.model";

@Component(
      {
         templateUrl: './host.summary.component.html',
         styleUrls: ['./host.summary.component.scss']
      }
)

export class HostSummaryComponent implements OnInit {
   private static readonly HOST_SUMMARY_VIEW_NAVIGATION_ID = "hostMonitor";
   private loading: boolean = true;
   private contextObjectId: string;
   numberOfRelaterChassis: number;

   constructor(@Inject(ChassisService) private chassisService: ChassisService) {
   }

   ngOnInit(): void {
      this.contextObjectId =
            this.chassisService.htmlClientSdk.app.getContextObjects()[0].id;
      this.loadData();
   }

   public navigateToHostMonitorView(): void {
      const navigateParams: any = {
         targetViewId: HostSummaryComponent.HOST_SUMMARY_VIEW_NAVIGATION_ID,
         objectId: this.contextObjectId
      };
      this.chassisService.htmlClientSdk.app.navigateTo(navigateParams);
   }

   private loadData(): void {
      this.chassisService.getRelatedChassis(this.contextObjectId)
            .then((chassis: Chassis[]) => {
               this.numberOfRelaterChassis = chassis.length;
               this.loading = false;
            })
   }
}
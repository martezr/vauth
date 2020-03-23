/* Copyright (c) 2018 VMware, Inc. All rights reserved. */
import { Component, Inject, OnInit } from "@angular/core";
import { Chassis } from "../../model/chassis.model";
import { ChassisService } from "../../services/chassis.service";
import { HostsService } from "../../services/hosts.service";
import { Host } from "../../model/host.model";
import 'rxjs/add/operator/finally';

@Component(
      {
         templateUrl: './host.monitor.component.html'
      }
)
export class HostMonitorComponent implements OnInit {

   private chassisList: Chassis[];
   private _selectedChassis: Chassis[];
   private loading: boolean = true;
   private contextObjectId: string;

   constructor(@Inject(ChassisService) private chassisService: ChassisService,
               @Inject(HostsService) private hostsService: HostsService) {
   }

   ngOnInit(): void {
      this.contextObjectId =
            this.chassisService.htmlClientSdk.app.getContextObjects()[0].id;
      this.loadData();
   }

   get selectedChassis() {
      return this._selectedChassis;
   }

   set selectedChassis(chassis: Chassis[]) {
      this._selectedChassis = chassis;
   }

   public updateChassisRelation() {
      this.loading = true;
      const host: Host = new Host();
      host.id = this.contextObjectId;
      host.relatedChassisIds = this.selectedChassis.map((chassis: Chassis) => {
         return chassis.id;
      });
      this.hostsService.edit(host).then(
            (resolve) => this.loading = false,
            (reject) => this.loading = false);
   }

   private loadData(): void {
      this.chassisService.getAllChassis().then((chassis: Chassis[]) => {
         this.chassisList = chassis;
         this.selectedChassis = this.filterRelatedChassis();
         this.loading = false;
      });
   }

   private filterRelatedChassis(): Chassis[] | null {
      if (!this.chassisList || this.chassisList.length < 1) {
         return null;
      }
      return this.chassisList.filter((chassis: Chassis) => {
         return !!chassis.relatedHostsIds &&
               chassis.relatedHostsIds.indexOf(this.contextObjectId) > -1;
      });
   }
}
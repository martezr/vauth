/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {
   Component, EventEmitter, Input, OnInit, Output
} from '@angular/core';

import { HostsService } from "../../../services/hosts.service";
import { Host } from "../../../model/host.model";
import { Chassis } from "../../../model/chassis.model";

@Component(
      {
         selector: 'hosts-list-view',
         templateUrl: './hosts-list.component.html'
      }
)

export class HostListComponent implements OnInit {

   private _selectedHosts: Host[] = <Host[]>[];

   @Input()
   preselectedHostsIds: string[];

   @Input()
   chassis: Chassis;

   @Output()
   hostsSelectionChange = new EventEmitter<string[]>();

   @Output()
   onNavigateToHostObject = new EventEmitter<any>();

   loading: boolean = false;
   connectedHosts: Host[];
   onContextMenu = HostListComponent.preventContextMenu;

   constructor(private hostsService: HostsService) {
   }

   ngOnInit(): void {
      this.retrieveHosts();
   }

   /**
    * Setter of the two-way binding with the Datagrid's selected items
    * @param selectedHosts - array of the updated Datagrid's selection
    */
   set selectedHosts(selectedHosts: Host[]) {
      this._selectedHosts = selectedHosts;
      if (!!selectedHosts) {
         this.emitHostSelectionChangeEvent(selectedHosts);
      }
   }

   /**
    * Getter of the two-way binding with the Datagrid's selected items
    */
   get selectedHosts(): Host[] {
      return this._selectedHosts;
   }

   /**
    * Navigate To the host summary view of a given objectId
    */
   navigateToHostObject(objectId: string): void {
      let navigateParams = {
         objectId: objectId
      };
      this.hostsService.htmlClientSdk.app.navigateTo(navigateParams);
      this.onNavigateToHostObject.emit();
   }

   /**
    * Refresh the list of host objects.
    */
   private retrieveHosts(): void {
      this.loading = true;
      this.hostsService.getConnectedHosts(this.chassis)
            .then((result: Host[]) => {
               this.connectedHosts = result;
               this.selectedHosts =
                     this.filterPreselectedHosts(this.connectedHosts);
               this.loading = false;
            });
   }

   /**
    * Filter out an array of preselected Host objects out of all connected
    * Hosts objects
    * @param hostsList
    */
   private filterPreselectedHosts(hostsList: Host[]): Host[] {
      if (!this.preselectedHostsIds) {
         return null;
      }
      return hostsList.filter((host: Host) =>
            this.preselectedHostsIds.indexOf(host.id) >= 0);
   }

   /**
    * Notify the consumers that Host objects selection has changed.
    * @param selectedHosts
    */
   private emitHostSelectionChangeEvent(selectedHosts: Host[]) {
      this.hostsSelectionChange.emit(selectedHosts.map((host: Host) => {
         return host.id;
      }));
   }

   static preventContextMenu(): boolean {
      return false;
   }
}

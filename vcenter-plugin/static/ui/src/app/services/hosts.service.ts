/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import { Injectable } from "@angular/core";
import { Chassis } from "../model/chassis.model";
import { GlobalService } from "./global.service";

import "rxjs/add/operator/toPromise";
import { HttpService } from "./http.service";
import { Host } from "../model/host.model";

@Injectable()
export class HostsService extends GlobalService {
   constructor(private httpService: HttpService) {
      super();
   }

   /**
    * Sends a get message to get all connected hosts
    */
   public getConnectedHosts(chassis: Chassis): Promise<Host[]> {
      const endpoint = chassis ? `chassis/${chassis.id}/hosts` : "hosts";
      return new Promise((resolve, reject) => {
         this.httpService.get(endpoint)
               .subscribe((result: any) => {
                  if (result) {
                     resolve(result);
                  } else {
                     reject("Failed to get a list of hosts.");
                  }
               });
      });
   }

   /**
    * Sends a message to edit the Host object
    */
   public edit(host: Host): Promise<any> {
      const endpoint = "hosts";
      return new Promise((resolve, reject) => {
         this.httpService.put(endpoint, host)
               .subscribe((result: any) => {
                  if (!result) {
                     reject("Failed to assign selected chassis to host.");
                  } else {
                     resolve();
                  }
               });
      });
   }
}

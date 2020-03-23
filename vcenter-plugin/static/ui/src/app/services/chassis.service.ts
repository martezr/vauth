/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import { Injectable } from "@angular/core";
import { Chassis } from "../model/chassis.model";
import { GlobalService } from "./global.service";

import "rxjs/add/operator/toPromise";
import { HttpService } from "./http.service";
import { URLSearchParams } from "../../../node_modules/@angular/http/src/url_search_params";

@Injectable()
export class ChassisService extends GlobalService {

   constructor(private httpService: HttpService) {
      super();
   }

   /**
    * Creates a new object of type Chassis
    * @param chassis - the created object.
    */
   public create(chassis: Chassis): Promise<Chassis | null> {
      chassis.name = chassis.name.trim();
      return new Promise((resolve, reject) => {
         this.httpService.post("chassis/create", JSON.stringify(chassis))
               .subscribe((result: any) => {
                  if (result) {
                     chassis.id = result;
                     chassis.healthStatus = 45;
                     chassis.complianceStatus = 81;
                     resolve(chassis);
                  } else {
                     reject("Failed to create chassis.")
                  }
               })
      });
   }

   /**
    * Edit the given chassis.
    * @param chassis - the edited chassis.
    */
   public edit(chassis: Chassis): Promise<boolean> {
      let newChassis = Object.assign(new Chassis(), chassis);
      newChassis.name = newChassis.name.trim();
      return new Promise((resolve, reject) => {
         this.httpService.put("chassis/edit", JSON.stringify(chassis))
               .subscribe((result: any) => {
                  resolve(result);
               });
      });
   }

   public delete(): Promise<boolean> {
      const chassisIds: string =
            this.htmlClientSdk.app.getContextObjects().join(",");
      return new Promise((resolve, reject) => {
         this.httpService.delete("chassis/delete", `ids=${chassisIds}`)
               .subscribe(() => {
                  resolve(true);
               });
      });
   }

   /**
    * Retrieves all related Chassis to the provided objectId
    * @param objectId
    */
   public getRelatedChassis(objectId: string): Promise<Chassis[]> {
      return new Promise((resolve, reject) => {
         this.httpService.get(`hosts/${objectId}/chassis`)
               .subscribe((result: Chassis[]) => {
                  if (result) {
                     resolve(result);
                  } else {
                     reject(`Failed to get related chassis for object '${objectId}'.`);
                  }
               });
      });
   }

   public getAllChassis(): Promise<Chassis[]> {
      let data: Chassis[];

      return new Promise((resolve, reject) => {
         this.httpService.get("chassis/list")
               .subscribe((result: any) => {
                  if (result) {
                     data = result as Chassis[];
                     for (let chassis of data) {
                        chassis.healthStatus = 45;
                        chassis.complianceStatus = 81;
                     }
                     resolve(data);
                  } else {
                     reject("Failed to get a list of all chassis.")
                  }
               });
      });
   }
}


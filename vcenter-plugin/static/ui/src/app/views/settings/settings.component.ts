/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import { Component } from "@angular/core";
import { OnInit } from "@angular/core";
import { Chassis } from "../../model/chassis.model";

@Component(
      {
         templateUrl: "./settings.component.html"
      }
)

export class SettingsComponent implements OnInit {
   numberOfChassisPerPage: number =
         Chassis.DEFAULT_CHASSIS_PAGE_SIZE;
   showSuccessMessage: boolean = false;

   ngOnInit(): void {
      let value =
            parseInt(localStorage.getItem(Chassis.PROP_CHASSIS_PAGE_SIZE));
      if (value && value > 0) {
         this.numberOfChassisPerPage = value;
      }
   }

   /**
    * Sets the new value in the local storage.
    * @param numberChassisPerPage -
    * number of chassis displayed in the chassis list per page.
    */
   private setNumberChassisPerPageInLocalStorage(numberChassisPerPage: number) {
      localStorage.setItem(
            Chassis.PROP_CHASSIS_PAGE_SIZE, numberChassisPerPage.toString());
   }

   /**
    * Triggered when user clicks on "Update" button.
    */
   onUpdate() {
      let value: number = parseInt(this.numberOfChassisPerPage.toString());
      if (value === Number.NaN || value <= 0) {
         value = Chassis.DEFAULT_CHASSIS_PAGE_SIZE;
      }

      this.numberOfChassisPerPage = value;
      this.setNumberChassisPerPageInLocalStorage(value);
      this.showSuccessMessage = true;
   }
}

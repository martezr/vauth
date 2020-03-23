/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { HttpClientModule } from '@angular/common/http';
import { HttpClient } from "@angular/common/http";
import { ClarityModule } from 'clarity-angular';
import { ClarityIcons } from 'clarity-icons';

import { AppComponent } from './app.component';
import { ChassisService } from './services/chassis.service';
import { HostsService } from './services/hosts.service';
import { ModalConfigService } from './services/modal.service';
import { ResourceService } from './services/resource.service';
import { HeaderComponent } from './views/header/header.component';
import { StatusComponent } from './views/status/status.component';
import { ModalLabelInputComponent } from './views/modal-label-input/modal-label-input.component';
import { AppRoutingModule, routableComponents } from './app-routing.module';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { SummaryComponent } from "./views/tabs/summary/summary.component";
import { MonitorComponent } from "./views/tabs/monitor/monitor.component";
import { ConfigureComponent } from "./views/tabs/configure/configure.component";
import { HostListComponent } from "./views/tabs/hosts/hosts-list.component";
import { SettingsComponent } from "./views/tabs/configure/settings/settings.component";
import { AdditionalViewComponent } from "./views/tabs/configure/additional/additional-view.component";
import { TranslateModule, TranslateLoader } from "@ngx-translate/core";
import { TranslateHttpLoader } from '@ngx-translate/http-loader';

import 'clarity-icons/shapes/all-shapes';
import { WindowService } from "./services/window.service";
import { EntryPointComponent } from "./views/entry-point/entry.point.component";
import {HttpService} from "./services/http.service";
import {GlobalService} from "./services/global.service";
import { HostSummaryComponent } from "./views/host-views/host.summary.component";
import { HostMonitorComponent } from "./views/host-views/host.monitor.component";

export function createTranslateLoader(http: HttpClient) {
   return new TranslateHttpLoader(http, './assets/i18n/', '.json');
}

@NgModule({
   imports: [
      BrowserModule,
      ClarityModule,
      HttpClientModule,
      TranslateModule.forRoot({
         loader: {
            provide: TranslateLoader,
            useFactory: (createTranslateLoader),
            deps: [HttpClient]
         }
      }),
      FormsModule,
      HttpModule,
      AppRoutingModule,
      BrowserAnimationsModule
   ],
   declarations: [AppComponent,
      routableComponents,
      EntryPointComponent,
      SummaryComponent,
      MonitorComponent,
      ConfigureComponent,
      HostListComponent,
      SettingsComponent,
      AdditionalViewComponent,
      HeaderComponent,
      StatusComponent,
      ModalLabelInputComponent,
      HostSummaryComponent,
      HostMonitorComponent
   ],
   providers: [ChassisService, ModalConfigService, ResourceService, HostsService, WindowService, HttpService, GlobalService],
   bootstrap: [AppComponent]
})

export class AppModule {
   constructor() {
      ClarityIcons.add({
         "chassisIcon": `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><title>chassis</title><g id="083e8a78-1df0-436f-9a22-5827915fd5f3" data-name="Layer 1"><path d="M.53,9H14.47a.53.53,0,0,1,.53.53v2.93a.53.53,0,0,1-.53.53H.54A.54.54,0,0,1,0,12.47V9.54A.53.53,0,0,1,.53,9Z" fill="#cbd9e0"/><path d="M15,9,13.17,4.07c-.15-.29-.22-.51-.77-.51H3.6c-.54,0-.62.23-.77.52L1,9Z" fill="#cad8df"/><path d="M.54,9H14.47a.53.53,0,0,1,.53.53v2.93a.53.53,0,0,1-.53.53H.54A.54.54,0,0,1,0,12.47V9.55A.54.54,0,0,1,.54,9Z" fill="#9bb2ba"/><path d="M15.44,9H.66a1,1,0,0,1,1-1H14.48A1,1,0,0,1,15.44,9Z" fill="#fff"/><path d="M15,14H1a1,1,0,0,1-1-1V9.26L2,3.77A1.2,1.2,0,0,1,3.23,3h9.39a1.2,1.2,0,0,1,1.23.77l.06.12L16,9.56V13A1,1,0,0,1,15,14ZM3.26,4C3,4,3,4,2.93,4.22L1,9.69V13l14,0,0-3.42v0L13,4.23C12.89,4,12.89,4,12.66,4Z" fill="#235d80"/><rect x="2" y="10" width="1" height="2" fill="#235d80"/><rect x="4" y="10.04" width="0.91" height="1.93" fill="#235d80"/><rect x="6" y="10" width="1" height="2" fill="#235d80"/><rect x="8.02" y="10" width="3.98" height="1" fill="#235d80"/><rect x="13" y="10" width="1" height="1" fill="#235d80"/></g></svg>`
      });
   }
}

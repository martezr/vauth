/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Component, OnInit} from '@angular/core';
import {TranslateService} from '@ngx-translate/core';
import {ActivatedRoute, Params, Router} from "@angular/router";
import {GlobalService} from "./services/global.service";
import {Subscription} from "rxjs/Subscription";

declare const htmlClientSdk: any;

@Component({
   selector: 'app-root',
   templateUrl: './app.component.html',
   styleUrls: ['./app.component.scss']
})

export class AppComponent implements OnInit {
   private subscription: Subscription;

   private initialThemeLoadComplete: boolean = false;

   public get initialized(): boolean {
      return this.globalService.htmlClientSdk.isInitialized() &&
            this.initialThemeLoadComplete;
   }

   constructor(private translate: TranslateService,
               private router: Router,
               private route: ActivatedRoute,
               private globalService: GlobalService) {
   }

   ngOnInit(): void {
      this.translate.addLangs(["en-US", "de-DE", "fr-FR"]);
      this.translate.setDefaultLang('en-US');
      this.globalService.htmlClientSdk.initialize(() => {
         let locale = this.globalService.htmlClientSdk.app.getClientLocale();
         if (locale && this.translate.getLangs().indexOf(locale) > 0) {
            this.translate.use(locale);
         }

         this.subscription = this.route.queryParams.subscribe(
            (params: Params) => {
               let view = params['view'];
               if (view) {
                  // Replace the entry URL with the redirected one. The side effect is that
                  // in the browser's history will remain only one record. So one will be
                  // able to move back to the previous client's URL.
                  this.router.navigate(
                     ['/' + view, params],
                     { queryParams: params, replaceUrl: true });
               }
            });

         this.loadTheme(true, this.globalService.htmlClientSdk.app.getTheme());
         this.globalService.htmlClientSdk.event.onThemeChanged(
               this.loadTheme.bind(this, false)
         );
      });
   }

   private loadTheme(firstLoad: boolean, theme: any): void {
      let themeName: string = theme.name;
      let supportedThemeNames: string[] = ['light', 'dark'];
      if (supportedThemeNames.indexOf(themeName) === -1) {
         themeName = supportedThemeNames[0];
      }

      let styleSheetLinkElement =
            (<HTMLLinkElement> document.getElementById('theme-stylesheet-link'));
      let themeCssUrl = `theme-${themeName}.bundle.css`;

      if (firstLoad) {
         let initialThemeLoadCompleteListener = (event: Event) => {
            this.initialThemeLoadComplete = true;
            styleSheetLinkElement.removeEventListener('load', initialThemeLoadCompleteListener);
            styleSheetLinkElement.removeEventListener('error', initialThemeLoadCompleteListener);
         };

         styleSheetLinkElement.addEventListener('load', initialThemeLoadCompleteListener);
         styleSheetLinkElement.addEventListener('error', initialThemeLoadCompleteListener);
      }
      styleSheetLinkElement.setAttribute("href", themeCssUrl);

      document.documentElement.setAttribute("data-theme", themeName);
   }
}

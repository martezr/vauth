/* Copyright (c) 2018 VMware, Inc. All rights reserved. */

import {Component} from '@angular/core';
import {ChassisService} from "../../../services/chassis.service";

@Component(
    {
        selector: '[delete-dialog-content]',
        templateUrl: './delete.component.html'
    }
)

export class DeleteComponent {
    constructor(private chassisService: ChassisService) {
    }

    onSubmit(): void {
        this.chassisService.delete()
            .then((result: boolean) => {
                this.chassisService.htmlClientSdk.modal.close(result);
            })
            .catch(() => {
                this.chassisService.htmlClientSdk.modal.close();
            });
    }

    onCancel(): void {
        this.chassisService.htmlClientSdk.modal.close();
    }
}

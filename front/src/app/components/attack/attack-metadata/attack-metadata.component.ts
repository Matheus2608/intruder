import { Component } from '@angular/core';

@Component({
  selector: 'app-attack-metadata',
  imports: [],
  templateUrl: './attack-metadata.component.html',
  styleUrl: './attack-metadata.component.css'
})
export class AttackMetadataComponent {
    isError: any;
    requestId: any;
    payload: any;
    statusCode: any;
    timeElapsed: any;
    errorMsg: any;
    length: any;
}

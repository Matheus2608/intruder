import { Component, Input } from '@angular/core';

@Component({
  selector: '[app-attack-metadata]',
  imports: [],
  templateUrl: './attack-metadata.component.html',
  styleUrl: './attack-metadata.component.css'
})
export class AttackMetadataComponent {
    @Input() isError: any;
    @Input() requestId: any;
    @Input() payload: any;
    @Input() statusCode: any;
    @Input() timeElapsed: any;
    @Input() errorMsg: any;
    @Input() length: any;
}

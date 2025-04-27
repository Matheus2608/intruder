import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: '[app-payload-area]', // seletor de atributo
  standalone: true,
  template: '', // sem HTML aqui!
  styleUrls: ['./payload-area.component.css'],
  host: {
    '[attr.data-payload]': 'payloadId',
    '[name]': 'payloadName',
    'class': 'payload-input',
    '[placeholder]': '"Payload list " + payloadId',
    '[ngModel]': 'payloadText',
    '(ngModelChange)': 'payloadText = $event'
  }
})
export class PayloadAreaComponent {
  @Input() payloadId!: number;
  @Input() payloadName!: string;
  @Input() payloadText!: string;

  ngOnInit(): void {
    this.payloadName = "payload" + this.payloadId;
    console.log(this.payloadName);
  }
}


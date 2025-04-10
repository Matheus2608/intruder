import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'app-payload-area',
  imports: [],
  templateUrl: './payload-area.component.html',
  styleUrl: './payload-area.component.css'
})
export class PayloadAreaComponent implements OnInit {

  @Input() payloadId! : number;
  payloadName!: string;

  ngOnInit(): void {
    this.payloadName = "payload" + this.payloadId;
    console.log(this.payloadName);
  }

}

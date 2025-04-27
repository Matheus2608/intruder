import { Component, Input, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-payload-area',
  imports: [FormsModule],
  templateUrl: './payload-area.component.html',
  styleUrl: './payload-area.component.css'
})
export class PayloadAreaComponent implements OnInit {

  @Input() payloadId! : number;
  payloadName!: string;
  payloadText: string = "";

  ngOnInit(): void {
    this.payloadName = "payload" + this.payloadId;
    console.log(this.payloadName);
  }

}

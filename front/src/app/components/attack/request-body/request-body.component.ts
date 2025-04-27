import { Component } from '@angular/core';

@Component({
  selector: 'app-request-body',
  imports: [],
  templateUrl: './request-body.component.html',
  styleUrl: './request-body.component.css'
})
export class RequestBodyComponent {
  isError: any;
  errorMsg: any;
  httpReq: any;
  httpRes: any;

}

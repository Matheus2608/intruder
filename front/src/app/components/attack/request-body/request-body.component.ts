import { Component, Input } from '@angular/core';

@Component({
  selector: '[app-request-body]',
  imports: [],
  templateUrl: './request-body.component.html',
  styleUrl: './request-body.component.css'
})
export class RequestBodyComponent {
  @Input() isError: any;
  @Input() errorMsg: any;
  @Input() httpReq: any;
  @Input() httpRes: any;

}

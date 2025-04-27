import { Component, ElementRef, OnInit, Renderer2, ViewChild, ViewChildren } from '@angular/core';
import { PayloadAreaComponent } from '../payload-area/payload-area.component';
import { ApiClientService } from '../../services/api-client.service';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AttackInput } from '../../interfaces/AttackInput';
import { AttackOutput } from '../../interfaces/AttackOutput';

@Component({
  selector: 'app-index',
  imports: [FormsModule, PayloadAreaComponent],
  templateUrl: './index.component.html',
  styleUrl: './index.component.css'
})
export class IndexComponent implements OnInit {

  ngOnInit(): void {
    let initialHttpPayloadLines : string[] = [
      "Paste your raw HTTP request here!",
      'Select the parts you want to replace with payloads (like "admin" or "password" above) and click "Add §".',
      "Each marked section (§...§) corresponds to a payload list on the right.",
      "Ensure markers surround the exact value, e.g., `Cookie: session=§value§`, not `§Cookie:§ session=value`.",
      "",
      "Example: ",
      "POST /login HTTP/1.1",
      "Host: example.com",
      "Content-Type: application/x-www-form-urlencoded",
      "Content-Length: 27",
      "",
      "username=§admin§&password=§password§",
    ]

    this.initialHttpPayload = initialHttpPayloadLines.join("\n")
  }

  @ViewChild('httpPayload') httpPayload!: ElementRef;
  @ViewChildren('payloadArea') payloadAreas!: PayloadAreaComponent[];
  indexPayload : number[] = [];
  initialHttpPayload : string = "";
  selectedAttackType = "sniper";

  constructor(private apiClient : ApiClientService, private router : Router, private renderer: Renderer2) {}

  addSpecialCharsBetweenSelectedArea = () => {
    const selection = window.getSelection();
    if (selection && selection.rangeCount > 0) {
        const range = selection.getRangeAt(0);
        const selectedText = range.toString();
        const wrappedText = `§${selectedText}§`;

        const span = document.createElement('span');
        span.textContent = wrappedText;
        span.style.color='red'
        range.deleteContents();
        range.insertNode(span);

        selection.removeAllRanges(); // Remove a seleção
    }

    this.indexPayload.push(this.indexPayload.length + 1);
  }

  removeDollarSpans() {
      const div = this.httpPayload.nativeElement;
      const spans = div.querySelectorAll('span');
      console.log(spans);

      spans.forEach((span: HTMLSpanElement) => {
          const textContent = span.textContent;

          if (textContent?.startsWith('§') && textContent?.endsWith('§')) {
                const textNode = this.renderer.createText(textContent.slice(1, -1));
                const range = document.createRange();
                range.selectNode(span);
                range.deleteContents();
                range.insertNode(textNode);
          }
      });

      this.indexPayload = [];
  }

  updateNumberOfPayloads() {
    const text = this.httpPayload.nativeElement.innerText;
    const regex = /§/g;
    const numberOfSpecialCharacters = text.match(regex)?.length || 0;
    const numberOfPayloads = Math.floor(numberOfSpecialCharacters / 2);
    this.indexPayload = Array.from({ length: numberOfPayloads }, (_, i) => i + 1);
  }

  attack(event: Event) {
    event.preventDefault();

    this.apiClient.getFakeApi(this.createRequest()).subscribe({
      next: (response : AttackOutput) => {
        console.log('Response:', response);
        this.router.navigate(['/attack']);
      },
      error: (error) => {
        console.error('Error:', error);
      }
    });
  }

  createRequest() : AttackInput {
    return {
      typeOfAttack : this.selectedAttackType,
      path : "api/fake",
      payloads: this.getAllPayloads(),
      httpRequest: this.httpPayload.nativeElement.innerText,
    }
  }

  getAllPayloads() : string[] {
    return this.payloadAreas.map((payloadArea) => payloadArea.payloadText);
  }

}

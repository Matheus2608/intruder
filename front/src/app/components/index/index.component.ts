import { Component, ElementRef, OnInit, Renderer2, ViewChild } from '@angular/core';
import { PayloadAreaComponent } from '../payload-area/payload-area.component';

@Component({
  selector: 'app-index',
  imports: [PayloadAreaComponent],
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
  initialHttpPayload : string = "";

  constructor(private renderer: Renderer2) {}

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
  }

}

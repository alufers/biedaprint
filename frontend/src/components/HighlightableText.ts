import Vue, { CreateElement, VNode } from "vue";
import Component from "vue-class-component";
import { Inject, Watch } from "vue-property-decorator";

/**
 * HighlightableText is a vue component (it uses a render function so no .vue file is required) which renders a <span> and highlights the keywords provided by HighlightableTextZone in the context.
 * It only highlights the top-level text nodes in its default children slot.
 */
@Component<HighlightableText>({})
export default class HighlightableText extends Vue {
  @Inject("$highlightsComputed")
  $highlightsComputed: Function;

  get highlights(): string[] {
    return this.$highlightsComputed();
  }
  render(createElement: CreateElement) {
    return createElement(
      "span",
      this.processVNodeList(this.$slots.default, createElement)
    );
  }

  processVNodeList(vnodes: VNode[], createElement: CreateElement): VNode[] {
    return vnodes.map(v => {
      if (v.text) {
        let text = v.text;
        let textLower = v.text.toLowerCase();
        let highlightIndexes = [];
        let highlightLengths = [];
        for (let highlight of this.highlights) {
          if (highlight.length === 0) {
            continue;
          }
          let startIndex = 0;
          let index = 0;
          while (
            (index = textLower.indexOf(highlight.toLowerCase(), startIndex)) >
            -1
          ) {
            highlightIndexes.push(index);
            highlightLengths.push(highlight.length);
            startIndex = index + highlight.length;
          }
        }
        // create a highlight mask to avoid any overlap problems
        let highlightMask: boolean[] = Array(text.length);
        for (let i = 0; i < text.length; i++) {
          highlightMask[i] = false;
        }
        for (let i = 0; i < highlightIndexes.length; i++) {
          for (let j = 0; j < highlightLengths[i]; j++) {
            highlightMask[highlightIndexes[i] + j] = true;
          }
        }
        let resultingVnodes = [];
        let currNodeText = "";
        let currNodeHighlighted = false;

        for (let i = 0; i < highlightMask.length; i++) {
          currNodeHighlighted = highlightMask[i];
          currNodeText = text[i];
          while (true) {
            if (currNodeHighlighted !== highlightMask[i + 1]) {
              break;
            }
            i++;
            currNodeText += text[i];
          }

          if (currNodeHighlighted) {
            resultingVnodes.push(
              createElement(
                "span",
                {
                  class: "highlightable-text-highlighted"
                },
                [currNodeText]
              )
            );
          } else {
            resultingVnodes.push(currNodeText);
          }
        }
        return createElement("span", {}, resultingVnodes);
      }

      return v;
    });
  }
}

import type { ReactNode } from "react";

/**
 * Adds <wbr>s to URLs to allow them to word-wrap more logically when displayed.
 */
export function allowBreakingOnDots(stringWithDots: string | undefined): ReactNode[] {
  if (!stringWithDots) {
    return [];
  }
  const components = stringWithDots.split(".");
  const nodes: ReactNode[] = [];
  for (let i = 0; i < components.length; ++i) {
    nodes.push(components[i]);
    nodes.push(<wbr />);
    if (i !== components.length - 1) {
      nodes.push(".");
    }
  }
  return nodes;
}

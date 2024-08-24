import React from "react";
import styled from "styled-components";

// Styled Components
const RemarksContainer = styled.div`
  display: grid;
  grid-template-columns: ${({ grid }) => grid};
  color: #666;
  & > div {
    margin-right: 1em;
  }
`;

const RemarksText = styled.div`
  font-size: 0.9em;
`;

export const Remarks = ({ remarks, grid = "4em 1fr" }) =>
  remarks ? (
    <RemarksContainer grid={grid}>
      <div />
      <RemarksText>{remarks}</RemarksText>
    </RemarksContainer>
  ) : null;

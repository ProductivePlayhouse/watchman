import React from "react";
import styled from "styled-components";
import * as C from "Components";

export const HighestMatch = ({ data }) =>
{
  if (!data) return null;
  return (
    <C.Section>
      <C.SectionTitle>Highest Match</C.SectionTitle>
      {data === 1 && <h3>100% = Exact Match</h3>}
      {data < 1 && data >= 0.97 && <h3>{data * 100}% = Partial Match</h3>}
      {data < 0.97 && <h3>{data * 100}% = No Match</h3>}
    </C.Section>
  );
};

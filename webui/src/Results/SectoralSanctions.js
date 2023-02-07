import React from "react";
import { makeStyles } from "@material-ui/core/styles";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Paper from "@material-ui/core/Paper";

const useStyles = makeStyles({
  table: {
    minWidth: 650,
  },
});

export const SectoralSanctions = ({ data }) =>
{
  const classes = useStyles();

  if (!data) return null;

  return (
    <TableContainer component={Paper}>
      <Table className={classes.table} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>Entity ID</TableCell>
            <TableCell>Type</TableCell>
            <TableCell>Programs</TableCell>
            <TableCell>Name</TableCell>
            <TableCell>Addresses</TableCell>
            <TableCell>Remarks</TableCell>
            <TableCell>Alternate Names</TableCell>
            <TableCell>IDs</TableCell>
            <TableCell>Source List URL</TableCell>
            <TableCell>Source Info URL</TableCell>
            <TableCell>Match</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {data.map((s) => (
            <TableRow key={s.entityID}>
              <TableCell>{s.entityID}</TableCell>
              <TableCell>{s.Type}</TableCell>
              <TableCell>{s.Programs.join(", ")}</TableCell>
              <TableCell>{s.Name}</TableCell>
              <TableCell>{s.Addresses.join(", ")}</TableCell>
              <TableCell>{s.Remarks.join(", ")}</TableCell>
              <TableCell>{s.AlternateNames.join(", ")}</TableCell>
              <TableCell>{s.Ids.join(", ")}</TableCell>
              <TableCell>{s.SourceListURL}</TableCell>
              <TableCell>{s.SourceInfoURL}</TableCell>
              <TableCell>{s.Match}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};
// Copyright (c) 2020-2021 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	backend "github.com/decred/politeia/politeiad/backendv2"
	"github.com/decred/politeia/politeiad/backendv2/tstorebe/tstore"
	cmv1 "github.com/decred/politeia/politeiawww/api/comments/v1"
)

// CommentPolicy sends a comments v1 Policy request to politeiawww.
func (c *Client) CommentPolicy() (*cmv1.PolicyReply, error) {
	resBody, err := c.makeReq(http.MethodPost,
		cmv1.APIRoute, cmv1.RoutePolicy, nil)
	if err != nil {
		return nil, err
	}

	var pr cmv1.PolicyReply
	err = json.Unmarshal(resBody, &pr)
	if err != nil {
		return nil, err
	}

	return &pr, nil
}

// CommentNew sends a comments v1 New request to politeiawww.
func (c *Client) CommentNew(n cmv1.New) (*cmv1.NewReply, error) {
	resBody, err := c.makeReq(http.MethodPost,
		cmv1.APIRoute, cmv1.RouteNew, n)
	if err != nil {
		return nil, err
	}

	var nr cmv1.NewReply
	err = json.Unmarshal(resBody, &nr)
	if err != nil {
		return nil, err
	}

	return &nr, nil
}

// CommentVote sends a comments v1 Vote request to politeiawww.
func (c *Client) CommentVote(v cmv1.Vote) (*cmv1.VoteReply, error) {
	resBody, err := c.makeReq(http.MethodPost,
		cmv1.APIRoute, cmv1.RouteVote, v)
	if err != nil {
		return nil, err
	}

	var vr cmv1.VoteReply
	err = json.Unmarshal(resBody, &vr)
	if err != nil {
		return nil, err
	}

	return &vr, nil
}

// CommentDel sends a comments v1 Del request to politeiawww.
func (c *Client) CommentDel(d cmv1.Del) (*cmv1.DelReply, error) {
	resBody, err := c.makeReq(http.MethodPost,
		cmv1.APIRoute, cmv1.RouteDel, d)
	if err != nil {
		return nil, err
	}

	var dr cmv1.DelReply
	err = json.Unmarshal(resBody, &dr)
	if err != nil {
		return nil, err
	}

	return &dr, nil
}

// CommentCount sends a comments v1 Count request to politeiawww.
func (c *Client) CommentCount(cc cmv1.Count) (*cmv1.CountReply, error) {
	resBody, err := c.makeReq(http.MethodPost,
		cmv1.APIRoute, cmv1.RouteCount, cc)
	if err != nil {
		return nil, err
	}

	var cr cmv1.CountReply
	err = json.Unmarshal(resBody, &cr)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}

// Comments sends a comments v1 Comments request to politeiawww.
func (c *Client) Comments(cm cmv1.Comments) (*cmv1.CommentsReply, error) {
	resBody, err := c.makeReq(http.MethodPost,
		cmv1.APIRoute, cmv1.RouteComments, cm)
	if err != nil {
		return nil, err
	}

	var cr cmv1.CommentsReply
	err = json.Unmarshal(resBody, &cr)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}

// CommentVotes sends a comments v1 Votes request to politeiawww.
func (c *Client) CommentVotes(v cmv1.Votes) (*cmv1.VotesReply, error) {
	resBody, err := c.makeReq(http.MethodPost,
		cmv1.APIRoute, cmv1.RouteVotes, v)
	if err != nil {
		return nil, err
	}

	var vr cmv1.VotesReply
	err = json.Unmarshal(resBody, &vr)
	if err != nil {
		return nil, err
	}

	return &vr, nil
}

// CommentTimestamps sends a comments v1 Timestamps request to politeiawww.
func (c *Client) CommentTimestamps(t cmv1.Timestamps) (*cmv1.TimestampsReply, error) {
	resBody, err := c.makeReq(http.MethodPost,
		cmv1.APIRoute, cmv1.RouteTimestamps, t)
	if err != nil {
		return nil, err
	}

	var tr cmv1.TimestampsReply
	err = json.Unmarshal(resBody, &tr)
	if err != nil {
		return nil, err
	}

	return &tr, nil
}

// VerifyCommentTimestamps verifies that all timestamps in a comments v1
// TimestampsReply are valid.
func VerifyCommentTimestamps(tr cmv1.TimestampsReply) error {
	for cid, cts := range tr.Comments {
		// Verify comment adds
		for i, ts := range cts.Adds {
			err := tstore.VerifyTimestamp(convertCommentTimestamp(ts))
			if err != nil {
				return fmt.Errorf("verify comment add timestamp %v %v: %v",
					cid, i, err)
			}
		}

		// Verify comment del if one exists
		if cts.Del == nil {
			continue
		}
		err := tstore.VerifyTimestamp(convertCommentTimestamp(*cts.Del))
		if err != nil {
			return fmt.Errorf("verify comment del timestamp %v: %v",
				cid, err)
		}
	}
	return nil
}

func convertCommentProof(p cmv1.Proof) backend.Proof {
	return backend.Proof{
		Type:       p.Type,
		Digest:     p.Digest,
		MerkleRoot: p.MerkleRoot,
		MerklePath: p.MerklePath,
		ExtraData:  p.ExtraData,
	}
}

func convertCommentTimestamp(t cmv1.Timestamp) backend.Timestamp {
	proofs := make([]backend.Proof, 0, len(t.Proofs))
	for _, v := range t.Proofs {
		proofs = append(proofs, convertCommentProof(v))
	}
	return backend.Timestamp{
		Data:       t.Data,
		Digest:     t.Digest,
		TxID:       t.TxID,
		MerkleRoot: t.MerkleRoot,
		Proofs:     proofs,
	}
}
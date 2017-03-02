package storage

import (
	"path"

	"github.com/anacrolix/missinggo/resource"

	"github.com/anacrolix/torrent/metainfo"
)

type piecePerResource struct {
	p resource.Provider
}

func NewResourcePieces(p resource.Provider) ClientImpl {
	return &piecePerResource{
		p: p,
	}
}

func (s *piecePerResource) OpenTorrent(info *metainfo.Info, infoHash metainfo.Hash) (TorrentImpl, error) {
	return s, nil
}

func (s *piecePerResource) Close() error {
	return nil
}

func (piecePerResource) DeletePiece(chunks int) (int64, error) {
	return 0, nil
}

func (s *piecePerResource) Piece(p metainfo.Piece) PieceImpl {
	completed, err := s.p.NewInstance(path.Join("completed", p.Hash().HexString()))
	if err != nil {
		panic(err)
	}
	incomplete, err := s.p.NewInstance(path.Join("incomplete", p.Hash().HexString()))
	if err != nil {
		panic(err)
	}
	return piecePerResourcePiece{
		p: p,
		c: completed,
		i: incomplete,
	}
}

type piecePerResourcePiece struct {
	p metainfo.Piece
	c resource.Instance
	i resource.Instance
}

func (piecePerResourcePiece) DeletePiece(chunks int) (int64, error) {
	return 0, nil
}
func (s piecePerResourcePiece) GetIsComplete() bool {
	fi, err := s.c.Stat()
	return err == nil && fi.Size() == s.p.Length()
}

func (s piecePerResourcePiece) MarkComplete() error {
	return resource.Move(s.i, s.c)
}

func (s piecePerResourcePiece) MarkNotComplete() error {
	return s.c.Delete()
}

func (s piecePerResourcePiece) ReadAt(b []byte, off int64) (int, error) {
	if s.GetIsComplete() {
		return s.c.ReadAt(b, off)
	} else {
		return s.i.ReadAt(b, off)
	}
}

func (s piecePerResourcePiece) WriteAt(b []byte, off int64) (n int, err error) {
	return s.i.WriteAt(b, off)
}

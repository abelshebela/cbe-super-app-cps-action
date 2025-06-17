package domain

import "time"

type HQ struct {
  Name           string
  Enabled        bool
  IsDeleted      bool

  ArchiveExpiry  uint  
  BlockTime      uint  

  CreatedAt      time.Time
  LastModified   time.Time
}

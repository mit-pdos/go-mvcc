package tpcc

/* Warehouse table. */
type Warehouse struct {
	/* Primary key: W_ID */
	W_ID       uint8
	/* Data fields */
	W_NAME     [10]byte
	W_STREET_1 [20]byte
	W_STREET_2 [20]byte
	W_CITY     [20]byte
	W_STATE    [2]byte
	W_ZIP      [9]byte
	W_TAX      float32
	W_YTD      float32
}

/* District table. */
type District struct {
	/* Primary key: (D_W_ID, D_ID) */
	D_ID       uint8
	D_W_ID     uint8
	/* Data fields */
	D_NAME     [10]byte
	D_STREET_1 [20]byte
	D_STREET_2 [20]byte
	D_CITY     [20]byte
	D_STATE    [2]byte
	D_ZIP      [9]byte
	W_TAX      float32
	D_YTD      float32
}

/* Customer table. */
type Customer struct {
	/* Primary key: (C_ID, C_D_ID, C_W_ID) */
	C_ID          uint32
	C_D_ID        uint8
	C_W_ID        uint8
	/* TODO: data fields */
	C_LAST        [16]byte
	C_CREDIT      [2]byte
	C_BALANCE     float32
	C_YTD_PAYMENT float32
	C_PAYMENT_CNT uint16
	C_DATA        [500]byte
}

/* History table. */
type History struct {
	/* Primary key: H_ID (no primary key required in the spec) */
	H_ID     uint64
	/* Data fields */
	H_C_ID   uint32
	H_C_D_ID uint8
	H_C_W_ID uint8
	H_D_ID   uint8
	H_W_ID   uint8
	H_DATE   uint32
	H_AMOUNT float32
	H_DATA   [25]byte
}

/* NewOrder table. */
type NewOrder struct {
	/* Primary key: (NO_W_ID, NO_D_ID, NO_W_ID) */
	NO_O_ID uint32
	NO_D_ID uint8
	NO_W_ID uint8
	/* No data fields */
}

/* Order table. */
type Order struct {
	/* Primary key: (O_W_ID, O_D_ID, O_ID) */
	O_ID   uint32
	O_D_ID uint8
	O_W_ID uint8
	/* TODO: data fields */
}

/* OrderLine table. */
type OrderLine struct {
	/* Primary key: (OL_W_ID, OL_D_ID, OL_W_ID, OL_NUMBER) */
	OL_O_ID   uint32
	OL_D_ID   uint8
	OL_W_ID   uint8
	OL_NUMBER uint8
	/* TODO: data fields */
}

/* Item table. */
type Item struct {
	/* Primary key: I_ID */
	I_ID   uint32
	/* TODO: data fields */
}

/* Stock table. */
type Stock struct {
	/* Primary key: (S_W_ID, S_I_ID) */
	S_I_ID   uint32
	S_W_ID   uint8
	/* TODO: data fields */
}

const (
	TBLID_WAREHOUSE uint64 = iota << 56
	TBLID_DISTRICT
	TBLID_CUSTOMER
	TBLID_HISTORY
	TBLID_NEWORDER
	TBLID_ORDER
	TBLID_ORDERLINE
	TBLID_ITEM
	TBLID_STOCK
)
package epaper

import (
	"log"
	"time"

	rpio "github.com/stianeikeland/go-rpio/v4"
)

var (
	DIN_PIN  = rpio.Pin(10) // SPI0_MOSI
	CLK_PIN  = rpio.Pin(11) // SPIO_SCK
	RST_PIN  = rpio.Pin(17) // GPIO_17
	DC_PIN   = rpio.Pin(25) // GPIO_25
	CS_PIN   = rpio.Pin(8)  // GPIO_8
	BUSY_PIN = rpio.Pin(24) // GPIO_24
)

var (
	w = 0
	h = 0
)

func Init(width, height int) {
	w = width
	h = height
	if err := rpio.Open(); err != nil {
		panic(err)
	}
	RST_PIN.Output()
	DC_PIN.Output()
	CS_PIN.Output()
	BUSY_PIN.Input()

	// SPI bus = 0, device = 0
	rpio.SpiBegin(rpio.Spi0)
	rpio.SpiMode(0, 0)
	rpio.SpiSpeed(4000000)
	initModule()
}

func initModule() {
	Reset()
	WriteCmd(POWER_SETTING)
	WriteData(0x03, 0x00, 0x2b, 0x2b) // VDS_EN, VDG_EN | VCOM_HV, VGHL_LV[1], VGHL_LV[0] | VDH | VDL

	WriteCmd(BOOSTER_SOFT_START)
	WriteData(0x17, 0x17, 0x17)

	WriteCmd(POWER_ON)

	ReadBusy()
	WriteCmd(PANEL_SETTING)
	WriteData(0xbf) // KW-BF   KWR-AF  BWROTP 0f BWOTP 1f

	WriteCmd(PLL_CONTROL) // PLL setting
	WriteData(0x3c)       // 3A 100HZ   29 150Hz 39 200HZ  31 171HZ

	WriteCmd(RESOLUTION_SETTING) // resolution setting
	WriteData(0x01)              // 400
	WriteData(0x90)              // 128
	WriteData(0x01)              // 300
	WriteData(0x2c)

	WriteCmd(VCM_DC_SETTING) // vcom_DC setting
	WriteData(0x12)

	WriteCmd(0x50)  // VCOM AND DATA INTERVAL SETTING
	WriteData(0x97) // 97white border 77black border  VBDF 17|D7 VBDW 97 VBDB 57  VBDF F7 VBDW 77 VBDB 37  VBDR B7
	setLUT()
}

func setLUT() {
	WriteCmd(VCOM_LUT) // vcom
	WriteData(lut_vcom0...)

	WriteCmd(W2W_LUT) // ww --
	WriteData(lut_ww...)

	WriteCmd(B2W_LUT) // bw r
	WriteData(lut_bw...)

	WriteCmd(W2B_LUT) // wb w
	WriteData(lut_bb...)

	WriteCmd(B2B_LUT) // bb b
	WriteData(lut_wb...)
}

func SetPartialLut() {
	WriteCmd(VCOM_LUT)
	WriteData(EPD_4IN2_Partial_lut_vcom1...)

	WriteCmd(W2W_LUT)
	WriteData(EPD_4IN2_Partial_lut_ww1...)

	WriteCmd(B2W_LUT)
	WriteData(EPD_4IN2_Partial_lut_bw1...)

	WriteCmd(W2B_LUT)
	WriteData(EPD_4IN2_Partial_lut_wb1...)

	WriteCmd(B2B_LUT)
	WriteData(EPD_4IN2_Partial_lut_bb1...)
}

func SetGrayLUT() {
	WriteCmd(VCOM_LUT) // vcom
	WriteData(EPD_4IN2_4Gray_lut_vcom...)

	WriteCmd(W2W_LUT) // red not use
	WriteData(EPD_4IN2_4Gray_lut_ww...)

	WriteCmd(B2W_LUT) // bw r
	WriteData(EPD_4IN2_4Gray_lut_bw...)

	WriteCmd(W2B_LUT) // wb w
	WriteData(EPD_4IN2_4Gray_lut_wb...)

	WriteCmd(B2B_LUT) // bb b
	WriteData(EPD_4IN2_4Gray_lut_bb...)

	WriteCmd(0x25) // vcom
	WriteData(EPD_4IN2_4Gray_lut_ww...)
}

func Close() {
	rpio.SpiEnd(rpio.Spi0)
	RST_PIN.Low()
	DC_PIN.Low()
	rpio.Close()
}

func generateColorData(color uint8) []uint8 {
	var lineWidth = 0
	if w%8 == 0 {
		lineWidth = w / 8
	} else {
		lineWidth = w/8 + 1
	}
	var out = make([]uint8, 0)
	for i := 0; i < lineWidth*h; i++ {
		out = append(out, color)
	}
	return out
}

func Clear() {
	log.Println("Clear")
	WriteCmd(DATA_START_TRANSMISSION_1)
	WriteData(generateColorData(0xff)...)

	WriteCmd(DATA_START_TRANSMISSION_2)
	WriteData(generateColorData(0xff)...)

	WriteCmd(DISPLAY_REFRESH)
	ReadBusy()
}

func DrawData(imageData []byte) {
	lineWidth := w / 8
	if w%8 != 0 {
		lineWidth++
	}
	WriteCmd(B2B_LUT)
	for j := 0; j < h; j++ {
		for i := 0; i < lineWidth; i++ {
			WriteData(imageData[i+j*lineWidth])
		}
	}
	TurnOnDisplay()
}

func Reset() {
	log.Println("Reset")
	RST_PIN.High()
	time.Sleep(100 * time.Millisecond)
	RST_PIN.Low()
	time.Sleep(2 * time.Millisecond)
	RST_PIN.High()
	time.Sleep(100 * time.Millisecond)
}

func ReadBusy() {
	log.Println("ReadBusy", BUSY_PIN.Read())
	for BUSY_PIN.Read() == rpio.High {
		time.Sleep(10 * time.Millisecond)
	}
}

func TurnOnDisplay() {
	WriteCmd(DISPLAY_REFRESH)
	ReadBusy()
}

func WriteCmd(cmd byte) {
	DC_PIN.Low()
	CS_PIN.Low()
	rpio.SpiTransmit(cmd)
	CS_PIN.High()
}

func WriteData(data ...byte) {
	DC_PIN.High()
	CS_PIN.Low()
	rpio.SpiTransmit(data...)
	CS_PIN.High()
}

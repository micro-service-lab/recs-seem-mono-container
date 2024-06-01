package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/micro-service-lab/recs-seem-mono-container/app/entity"
	"github.com/micro-service-lab/recs-seem-mono-container/app/parameter"
	"github.com/micro-service-lab/recs-seem-mono-container/app/store"
)

// MimeTypeKey マイムタイプキー。
type MimeTypeKey string

const (
	// MimeTypeKeyUnknown 不明。
	MimeTypeKeyUnknown MimeTypeKey = "unknown"
	// MimeTypeKeyOctetStream application/octet-stream。
	MimeTypeKeyOctetStream MimeTypeKey = "application~octet-stream"
	// MimeTypeKeyXpixmap image/x-xpixmap。
	MimeTypeKeyXpixmap MimeTypeKey = "image~x-xpixmap"
	// MimeTypeKey7z application/x-7z-compressed。
	MimeTypeKey7z MimeTypeKey = "application~x-7z-compressed"
	// MimeTypeKeyZip application/zip。
	MimeTypeKeyZip MimeTypeKey = "application~zip"
	// MimeTypeKeyXlsx application/vnd.openxmlformats-officedocument.spreadsheetml.sheet。
	MimeTypeKeyXlsx MimeTypeKey = "application~vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	// MimeTypeKeyDocx application/vnd.openxmlformats-officedocument.wordprocessingml.document。
	MimeTypeKeyDocx MimeTypeKey = "application~vnd.openxmlformats-officedocument.wordprocessingml.document"
	// MimeTypeKeyPptx application/vnd.openxmlformats-officedocument.presentationml.presentation。
	MimeTypeKeyPptx MimeTypeKey = "application~vnd.openxmlformats-officedocument.presentationml.presentation"
	// MimeTypeKeyEpub application/epub+zip。
	MimeTypeKeyEpub MimeTypeKey = "application~epub~zip"
	// MimeTypeKeyJar application/jar。
	MimeTypeKeyJar MimeTypeKey = "application~jar"
	// MimeTypeKeyOdt application/vnd.oasis.opendocument.text。
	MimeTypeKeyOdt MimeTypeKey = "application~vnd.oasis.opendocument.text"
	// MimeTypeKeyOtt application/vnd.oasis.opendocument.text-template。
	MimeTypeKeyOtt MimeTypeKey = "application~vnd.oasis.opendocument.text-template"
	// MimeTypeKeyOds application/vnd.oasis.opendocument.spreadsheet。
	MimeTypeKeyOds MimeTypeKey = "application~vnd.oasis.opendocument.spreadsheet"
	// MimeTypeKeyOts application/vnd.oasis.opendocument.spreadsheet-template。
	MimeTypeKeyOts MimeTypeKey = "application~vnd.oasis.opendocument.spreadsheet-template"
	// MimeTypeKeyOdp application/vnd.oasis.opendocument.presentation。
	MimeTypeKeyOdp MimeTypeKey = "application~vnd.oasis.opendocument.presentation"
	// MimeTypeKeyOtp application/vnd.oasis.opendocument.presentation-template。
	MimeTypeKeyOtp MimeTypeKey = "application~vnd.oasis.opendocument.presentation-template"
	// MimeTypeKeyOdg application/vnd.oasis.opendocument.graphics。
	MimeTypeKeyOdg MimeTypeKey = "application~vnd.oasis.opendocument.graphics"
	// MimeTypeKeyOtg application/vnd.oasis.opendocument.graphics-template。
	MimeTypeKeyOtg MimeTypeKey = "application~vnd.oasis.opendocument.graphics-template"
	// MimeTypeKeyOdf application/vnd.oasis.opendocument.formula。
	MimeTypeKeyOdf MimeTypeKey = "application~vnd.oasis.opendocument.formula"
	// MimeTypeKeyOdc application/vnd.oasis.opendocument.chart。
	MimeTypeKeyOdc MimeTypeKey = "application~vnd.oasis.opendocument.chart"
	// MimeTypeKeySxc application/vnd.sun.xml.calc。
	MimeTypeKeySxc MimeTypeKey = "application~vnd.sun.xml.calc"
	// MimeTypeKeyPdf application/pdf。
	MimeTypeKeyPdf MimeTypeKey = "application~pdf"
	// MimeTypeKeyFdf application/vnd.fdf。
	MimeTypeKeyFdf MimeTypeKey = "application~vnd.fdf"
	// MimeTypeKeyOleStorage application/x-ole-storage。
	MimeTypeKeyOleStorage MimeTypeKey = "application~x-ole-storage"
	// MimeTypeKeyMsi application/x-ms-installer。
	MimeTypeKeyMsi MimeTypeKey = "application~x-ms-installer"
	// MimeTypeKeyAaf application/octet-stream。
	MimeTypeKeyAaf MimeTypeKey = "application~octet-stream~aaf"
	// MimeTypeKeyMsg application/vnd.ms-outlook。
	MimeTypeKeyMsg MimeTypeKey = "application~vnd.ms-outlook"
	// MimeTypeKeyXls application/vnd.ms-excel。
	MimeTypeKeyXls MimeTypeKey = "application~vnd.ms-excel"
	// MimeTypeKeyPub application/vnd.ms-publisher。
	MimeTypeKeyPub MimeTypeKey = "application~vnd.ms-publisher"
	// MimeTypeKeyPpt application/vnd.ms-powerpoint。
	MimeTypeKeyPpt MimeTypeKey = "application~vnd.ms-powerpoint"
	// MimeTypeKeyDoc application/msword。
	MimeTypeKeyDoc MimeTypeKey = "application~msword"
	// MimeTypeKeyPs application/postscript。
	MimeTypeKeyPs MimeTypeKey = "application~postscript"
	// MimeTypeKeyPsd image/vnd.adobe.photoshop。
	MimeTypeKeyPsd MimeTypeKey = "image~vnd.adobe.photoshop"
	// MimeTypeKeyP7s application/pkcs7-signature。
	MimeTypeKeyP7s MimeTypeKey = "application~pkcs7-signature"
	// MimeTypeKeyOgg application/ogg。
	MimeTypeKeyOgg MimeTypeKey = "application~ogg"
	// MimeTypeKeyOga audio/ogg。
	MimeTypeKeyOga MimeTypeKey = "audio~ogg"
	// MimeTypeKeyOgv video/ogg。
	MimeTypeKeyOgv MimeTypeKey = "video~ogg"
	// MimeTypeKeyPng image/png。
	MimeTypeKeyPng MimeTypeKey = "image~png"
	// MimeTypeKeyPngApng image/vnd.mozilla.apng。
	MimeTypeKeyPngApng MimeTypeKey = "image~vnd.mozilla.apng"
	// MimeTypeKeyJpg image/jpeg。
	MimeTypeKeyJpg MimeTypeKey = "image~jpeg"
	// MimeTypeKeyJxl image/jxl。
	MimeTypeKeyJxl MimeTypeKey = "image~jxl"
	// MimeTypeKeyJp2 image/jp2。
	MimeTypeKeyJp2 MimeTypeKey = "image~jp2"
	// MimeTypeKeyJpf image/jpx。
	MimeTypeKeyJpf MimeTypeKey = "image~jpx"
	// MimeTypeKeyJpm image/jpm。
	MimeTypeKeyJpm MimeTypeKey = "image~jpm"
	// MimeTypeKeyJxs image/jxs。
	MimeTypeKeyJxs MimeTypeKey = "image~jxs"
	// MimeTypeKeyGif image/gif。
	MimeTypeKeyGif MimeTypeKey = "image~gif"
	// MimeTypeKeyWebp image/webp。
	MimeTypeKeyWebp MimeTypeKey = "image~webp"
	// MimeTypeKeyExe application/vnd.microsoft.portable-executable。
	MimeTypeKeyExe MimeTypeKey = "application~vnd.microsoft.portable-executable"
	// MimeTypeKeyElf application/x-elf。
	MimeTypeKeyElf MimeTypeKey = "application~x-elf"
	// MimeTypeKeyObject application/x-object。
	MimeTypeKeyObject MimeTypeKey = "application~x-object"
	// MimeTypeKeyExecutable application/x-executable。
	MimeTypeKeyExecutable MimeTypeKey = "application~x-executable"
	// MimeTypeKeySharedlib application/x-sharedlib。
	MimeTypeKeySharedlib MimeTypeKey = "application~x-sharedlib"
	// MimeTypeKeyCoredump application/x-coredump。
	MimeTypeKeyCoredump MimeTypeKey = "application~x-coredump"
	// MimeTypeKeyArchive application/x-archive。
	MimeTypeKeyArchive MimeTypeKey = "application~x-archive"
	// MimeTypeKeyDeb application/vnd.debian.binary-package。
	MimeTypeKeyDeb MimeTypeKey = "application~vnd.debian.binary-package"
	// MimeTypeKeyTar application/x-tar。
	MimeTypeKeyTar MimeTypeKey = "application~x-tar"
	// MimeTypeKeyXar application/x-xar。
	MimeTypeKeyXar MimeTypeKey = "application~x-xar"
	// MimeTypeKeyBz2 application/x-bzip2。
	MimeTypeKeyBz2 MimeTypeKey = "application~x-bzip2"
	// MimeTypeKeyFits application/fits。
	MimeTypeKeyFits MimeTypeKey = "application~fits"
	// MimeTypeKeyTiff image/tiff。
	MimeTypeKeyTiff MimeTypeKey = "image~tiff"
	// MimeTypeKeyBmp image/bmp。
	MimeTypeKeyBmp MimeTypeKey = "image~bmp"
	// MimeTypeKeyIcon image/x-icon。
	MimeTypeKeyIcon MimeTypeKey = "image~x-icon"
	// MimeTypeKeyMpeg audio/mpeg。
	MimeTypeKeyMpeg MimeTypeKey = "audio~mpeg"
	// MimeTypeKeyFlac audio/flac。
	MimeTypeKeyFlac MimeTypeKey = "audio~flac"
	// MimeTypeKeyMidi audio/midi。
	MimeTypeKeyMidi MimeTypeKey = "audio~midi"
	// MimeTypeKeyApe audio/ape。
	MimeTypeKeyApe MimeTypeKey = "audio~ape"
	// MimeTypeKeyMpc audio/musepack。
	MimeTypeKeyMpc MimeTypeKey = "audio~musepack"
	// MimeTypeKeyAmr audio/amr。
	MimeTypeKeyAmr MimeTypeKey = "audio~amr"
	// MimeTypeKeyWav audio/wav。
	MimeTypeKeyWav MimeTypeKey = "audio~wav"
	// MimeTypeKeyAiff audio/aiff。
	MimeTypeKeyAiff MimeTypeKey = "audio~aiff"
	// MimeTypeKeyAu audio/basic。
	MimeTypeKeyAu MimeTypeKey = "audio~basic"
	// MimeTypeKeyMpegVideo video/mpeg。
	MimeTypeKeyMpegVideo MimeTypeKey = "video~mpeg"
	// MimeTypeKeyMov video/quicktime。
	MimeTypeKeyMov MimeTypeKey = "video~quicktime;mov"
	// MimeTypeKeyMqv video/quicktime。
	MimeTypeKeyMqv MimeTypeKey = "video~quicktime;mqv"
	// MimeTypeKeyMp4 video/mp4。
	MimeTypeKeyMp4 MimeTypeKey = "video~mp4"
	// MimeTypeKeyWebm video/webm。
	MimeTypeKeyWebm MimeTypeKey = "video~webm"
	// MimeTypeKey3gp video/3gpp。
	MimeTypeKey3gp MimeTypeKey = "video~3gpp"
	// MimeTypeKey3g2 video/3gpp2。
	MimeTypeKey3g2 MimeTypeKey = "video~3gpp2"
	// MimeTypeKeyAvi video/x-msvideo。
	MimeTypeKeyAvi MimeTypeKey = "video~x-msvideo"
	// MimeTypeKeyFlv video/x-flv。
	MimeTypeKeyFlv MimeTypeKey = "video~x-flv"
	// MimeTypeKeyMkv video/x-matroska。
	MimeTypeKeyMkv MimeTypeKey = "video~x-matroska"
	// MimeTypeKeyAsf video/x-ms-asf。
	MimeTypeKeyAsf MimeTypeKey = "video~x-ms-asf"
	// MimeTypeKeyAac audio/aac。
	MimeTypeKeyAac MimeTypeKey = "audio~aac"
	// MimeTypeKeyVoc audio/x-unknown。
	MimeTypeKeyVoc MimeTypeKey = "audio~x-unknown"
	// MimeTypeKeyMp4Audio audio/mp4。
	MimeTypeKeyMp4Audio MimeTypeKey = "audio~mp4"
	// MimeTypeKeyM4a audio/x-m4a。
	MimeTypeKeyM4a MimeTypeKey = "audio~x-m4a"
	// MimeTypeKeyM3u application/vnd.apple.mpegurl。
	MimeTypeKeyM3u MimeTypeKey = "application~vnd.apple.mpegurl"
	// MimeTypeKeyM4v video/x-m4v。
	MimeTypeKeyM4v MimeTypeKey = "video~x-m4v"
	// MimeTypeKeyRmvb application/vnd.rn-realmedia-vbr。
	MimeTypeKeyRmvb MimeTypeKey = "application~vnd.rn-realmedia-vbr"
	// MimeTypeKeyGz application/gzip。
	MimeTypeKeyGz MimeTypeKey = "application~gzip"
	// MimeTypeKeyClass application/x-java-applet。
	MimeTypeKeyClass MimeTypeKey = "application~x-java-applet"
	// MimeTypeKeySwf application/x-shockwave-flash。
	MimeTypeKeySwf MimeTypeKey = "application~x-shockwave-flash"
	// MimeTypeKeyCrx application/x-chrome-extension。
	MimeTypeKeyCrx MimeTypeKey = "application~x-chrome-extension"
	// MimeTypeKeyTtf font/ttf。
	MimeTypeKeyTtf MimeTypeKey = "font~ttf"
	// MimeTypeKeyWoff font/woff。
	MimeTypeKeyWoff MimeTypeKey = "font~woff"
	// MimeTypeKeyWoff2 font/woff2。
	MimeTypeKeyWoff2 MimeTypeKey = "font~woff2"
	// MimeTypeKeyOtf font/otf。
	MimeTypeKeyOtf MimeTypeKey = "font~otf"
	// MimeTypeKeyTtc font/collection。
	MimeTypeKeyTtc MimeTypeKey = "font~collection"
	// MimeTypeKeyEot application/vnd.ms-fontobject。
	MimeTypeKeyEot MimeTypeKey = "application~vnd.ms-fontobject"
	// MimeTypeKeyWasm application/wasm。
	MimeTypeKeyWasm MimeTypeKey = "application~wasm"
	// MimeTypeKeyShx application/vnd.shx。
	MimeTypeKeyShx MimeTypeKey = "application~vnd.shx"
	// MimeTypeKeyShp application/vnd.shp。
	MimeTypeKeyShp MimeTypeKey = "application~vnd.shp"
	// MimeTypeKeyDbf application/x-dbf。
	MimeTypeKeyDbf MimeTypeKey = "application~x-dbf"
	// MimeTypeKeyDcm application/dicom。
	MimeTypeKeyDcm MimeTypeKey = "application~dicom"
	// MimeTypeKeyRar application/x-rar-compressed。
	MimeTypeKeyRar MimeTypeKey = "application~x-rar-compressed"
	// MimeTypeKeyDjvu image/vnd.djvu。
	MimeTypeKeyDjvu MimeTypeKey = "image~vnd.djvu"
	// MimeTypeKeyMobi application/x-mobipocket-ebook。
	MimeTypeKeyMobi MimeTypeKey = "application~x-mobipocket-ebook"
	// MimeTypeKeyLit application/x-ms-reader。
	MimeTypeKeyLit MimeTypeKey = "application~x-ms-reader"
	// MimeTypeKeyBpg image/bpg。
	MimeTypeKeyBpg MimeTypeKey = "image~bpg"
	// MimeTypeKeySqlite application/vnd.sqlite3。
	MimeTypeKeySqlite MimeTypeKey = "application~vnd.sqlite3"
	// MimeTypeKeyDwg image/vnd.dwg。
	MimeTypeKeyDwg MimeTypeKey = "image~vnd.dwg"
	// MimeTypeKeyNes application/vnd.nintendo.snes.rom。
	MimeTypeKeyNes MimeTypeKey = "application~vnd.nintendo.snes.rom"
	// MimeTypeKeyLnk application/x-ms-shortcut。
	MimeTypeKeyLnk MimeTypeKey = "application~x-ms-shortcut"
	// MimeTypeKeyMacho application/x-mach-binary。
	MimeTypeKeyMacho MimeTypeKey = "application~x-mach-binary"
	// MimeTypeKeyQcp audio/qcelp。
	MimeTypeKeyQcp MimeTypeKey = "audio~qcelp"
	// MimeTypeKeyIcns image/x-icns。
	MimeTypeKeyIcns MimeTypeKey = "image~x-icns"
	// MimeTypeKeyHeic image/heic。
	MimeTypeKeyHeic MimeTypeKey = "image~heic"
	// MimeTypeKeyHeicSequence image/heic-sequence。
	MimeTypeKeyHeicSequence MimeTypeKey = "image~heic-sequence"
	// MimeTypeKeyHeif image/heif。
	MimeTypeKeyHeif MimeTypeKey = "image~heif"
	// MimeTypeKeyHeifSequence image/heif-sequence。
	MimeTypeKeyHeifSequence MimeTypeKey = "image~heif-sequence"
	// MimeTypeKeyHdr image/vnd.radiance。
	MimeTypeKeyHdr MimeTypeKey = "image~vnd.radiance"
	// MimeTypeKeyMrc application/marc。
	MimeTypeKeyMrc MimeTypeKey = "application~marc"
	// MimeTypeKeyMdb application/x-msaccess。
	MimeTypeKeyMdb MimeTypeKey = "application~x-msaccess~mdb"
	// MimeTypeKeyAccdb application/x-msaccess。
	MimeTypeKeyAccdb MimeTypeKey = "application~x-msaccess~accdb"
	// MimeTypeKeyZst application/zstd。
	MimeTypeKeyZst MimeTypeKey = "application~zstd"
	// MimeTypeKeyCab application/vnd.ms-cab-compressed。
	MimeTypeKeyCab MimeTypeKey = "application~vnd.ms-cab-compressed"
	// MimeTypeKeyRpm application/x-rpm。
	MimeTypeKeyRpm MimeTypeKey = "application~x-rpm"
	// MimeTypeKeyXz application/x-xz。
	MimeTypeKeyXz MimeTypeKey = "application~x-xz"
	// MimeTypeKeyLz application/lzip。
	MimeTypeKeyLz MimeTypeKey = "application~lzip"
	// MimeTypeKeyTorrent application/x-bittorrent。
	MimeTypeKeyTorrent MimeTypeKey = "application~x-bittorrent"
	// MimeTypeKeyCpio application/x-cpio。
	MimeTypeKeyCpio MimeTypeKey = "application~x-cpio"
	// MimeTypeKeyTzif application/tzif。
	MimeTypeKeyTzif MimeTypeKey = "application~tzif"
	// MimeTypeKeyXcf image/x-xcf。
	MimeTypeKeyXcf MimeTypeKey = "image~x-xcf"
	// MimeTypeKeyPat image/x-gimp-pat。
	MimeTypeKeyPat MimeTypeKey = "image~x-gimp-pat"
	// MimeTypeKeyGbr image/x-gimp-gbr。
	MimeTypeKeyGbr MimeTypeKey = "image~x-gimp-gbr"
	// MimeTypeKeyGlb model/gltf-binary。
	MimeTypeKeyGlb MimeTypeKey = "model~gltf-binary"
	// MimeTypeKeyAvif image/avif。
	MimeTypeKeyAvif MimeTypeKey = "image~avif"
	// MimeTypeKeyCabInstallshield application/x-installshield。
	MimeTypeKeyCabInstallshield MimeTypeKey = "application~x-installshield"
	// MimeTypeKeyJxr image/jxr。
	MimeTypeKeyJxr MimeTypeKey = "image~jxr"
	// MimeTypeKeyTxt text/plain。
	MimeTypeKeyTxt MimeTypeKey = "text~plain"
	// MimeTypeKeyHTML text/html。
	MimeTypeKeyHTML MimeTypeKey = "text~html"
	// MimeTypeKeySvg image/svg+xml。
	MimeTypeKeySvg MimeTypeKey = "image~svg+xml"
	// MimeTypeKeyXML text/xml。
	MimeTypeKeyXML MimeTypeKey = "text~xml"
	// MimeTypeKeyRss application/rss+xml。
	MimeTypeKeyRss MimeTypeKey = "application~rss~xml"
	// MimeTypeKeyAtom application/atom+xml。
	MimeTypeKeyAtom MimeTypeKey = "application~atom~xml"
	// MimeTypeKeyX3d model/x3d+xml。
	MimeTypeKeyX3d MimeTypeKey = "model~x3d+xml"
	// MimeTypeKeyKml application/vnd.google-earth.kml+xml。
	MimeTypeKeyKml MimeTypeKey = "application~vnd.google-earth.kml~xml"
	// MimeTypeKeyXlf application/x-xliff+xml。
	MimeTypeKeyXlf MimeTypeKey = "application~x-xliff~xml"
	// MimeTypeKeyDae model/vnd.collada+xml。
	MimeTypeKeyDae MimeTypeKey = "model~vnd.collada~xml"
	// MimeTypeKeyGml application/gml+xml。
	MimeTypeKeyGml MimeTypeKey = "application~gml~xml"
	// MimeTypeKeyGpx application/gpx+xml。
	MimeTypeKeyGpx MimeTypeKey = "application~gpx~xml"
	// MimeTypeKeyTcx application/vnd.garmin.tcx+xml。
	MimeTypeKeyTcx MimeTypeKey = "application~vnd.garmin.tcx~xml"
	// MimeTypeKeyAmf application/x-amf。
	MimeTypeKeyAmf MimeTypeKey = "application~x-amf"
	// MimeTypeKey3mf application/vnd.ms-package.3dmanufacturing-3dmodel+xml。
	MimeTypeKey3mf MimeTypeKey = "application~vnd.ms-package.3dmanufacturing-3dmodel~xml"
	// MimeTypeKeyXfdf application/vnd.adobe.xfdf。
	MimeTypeKeyXfdf MimeTypeKey = "application~vnd.adobe.xfdf"
	// MimeTypeKeyOwl application/owl+xml。
	MimeTypeKeyOwl MimeTypeKey = "application~owl~xml"
	// MimeTypeKeyPhp text/x-php。
	MimeTypeKeyPhp MimeTypeKey = "text~x-php"
	// MimeTypeKeyJs application/javascript。
	MimeTypeKeyJs MimeTypeKey = "application~javascript"
	// MimeTypeKeyLua text/x-lua。
	MimeTypeKeyLua MimeTypeKey = "text~x-lua"
	// MimeTypeKeyPl text/x-perl。
	MimeTypeKeyPl MimeTypeKey = "text~x-perl"
	// MimeTypeKeyPy text/x-python。
	MimeTypeKeyPy MimeTypeKey = "text~x-python"
	// MimeTypeKeyJSON application/json。
	MimeTypeKeyJSON MimeTypeKey = "application~json"
	// MimeTypeKeyGeojson application/geo+json。
	MimeTypeKeyGeojson MimeTypeKey = "application~geo~json"
	// MimeTypeKeyHar application/json。
	MimeTypeKeyHar MimeTypeKey = "application~json~har"
	// MimeTypeKeyNdjson application/x-ndjson。
	MimeTypeKeyNdjson MimeTypeKey = "application~x-ndjson"
	// MimeTypeKeyRtf text/rtf。
	MimeTypeKeyRtf MimeTypeKey = "text~rtf"
	// MimeTypeKeySrt application/x-subrip。
	MimeTypeKeySrt MimeTypeKey = "application~x-subrip"
	// MimeTypeKeyTcl text/x-tcl。
	MimeTypeKeyTcl MimeTypeKey = "text~x-tcl"
	// MimeTypeKeyCsv text/csv。
	MimeTypeKeyCsv MimeTypeKey = "text~csv"
	// MimeTypeKeyTsv text/tab-separated-values。
	MimeTypeKeyTsv MimeTypeKey = "text~tab-separated-values"
	// MimeTypeKeyVcf text/vcard。
	MimeTypeKeyVcf MimeTypeKey = "text~vcard"
	// MimeTypeKeyIcs text/calendar。
	MimeTypeKeyIcs MimeTypeKey = "text~calendar"
	// MimeTypeKeyWarc application/warc。
	MimeTypeKeyWarc MimeTypeKey = "application~warc"
	// MimeTypeKeyVtt text/vtt。
	MimeTypeKeyVtt MimeTypeKey = "text~vtt"
)

// MimeType マイムタイプ。
type MimeType struct {
	Key  string
	Name string
	Kind string
}

// MimeTypes マイムタイプ一覧。
//
//nolint:lll
var MimeTypes = []MimeType{
	{Key: string(MimeTypeKeyUnknown), Name: "未明", Kind: "unknown"},
	{Key: string(MimeTypeKeyOctetStream), Name: "Octet Stream", Kind: "application/octet-stream"},
	{Key: string(MimeTypeKeyXpixmap), Name: "QPixmap", Kind: "image/x-xpixmap"},
	{Key: string(MimeTypeKey7z), Name: "7z", Kind: "application/x-7z-compressed"},
	{Key: string(MimeTypeKeyZip), Name: "Zip", Kind: "application/zip"},
	{Key: string(MimeTypeKeyXlsx), Name: "Excel(xlsx)", Kind: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
	{Key: string(MimeTypeKeyDocx), Name: "Word(docx)", Kind: "application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
	{Key: string(MimeTypeKeyPptx), Name: "Power Point(pptx)", Kind: "application/vnd.openxmlformats-officedocument.presentationml.presentation"},
	{Key: string(MimeTypeKeyEpub), Name: "epub", Kind: "application/epub+zip"},
	{Key: string(MimeTypeKeyJar), Name: "jar", Kind: "application/jar"},
	{Key: string(MimeTypeKeyOdt), Name: "odt(OpenDocument Text)", Kind: "application/vnd.oasis.opendocument.text"},
	{Key: string(MimeTypeKeyOtt), Name: "ott(OpenDocument Text Template)", Kind: "application/vnd.oasis.opendocument.text-template"},
	{Key: string(MimeTypeKeyOds), Name: "ods(OpenDocument Spreadsheet)", Kind: "application/vnd.oasis.opendocument.spreadsheet"},
	{Key: string(MimeTypeKeyOts), Name: "ots(OpenDocument Spreadsheet Template)", Kind: "application/vnd.oasis.opendocument.spreadsheet-template"},
	{Key: string(MimeTypeKeyOdp), Name: "odp(OpenDocument Presentation)", Kind: "application/vnd.oasis.opendocument.presentation"},
	{Key: string(MimeTypeKeyOtp), Name: "otp(OpenDocument Presentation Template)", Kind: "application/vnd.oasis.opendocument.presentation-template"},
	{Key: string(MimeTypeKeyOdg), Name: "odg(OpenDocument Graphics)", Kind: "application/vnd.oasis.opendocument.graphics"},
	{Key: string(MimeTypeKeyOtg), Name: "otg(OpenDocument Graphics Template)", Kind: "application/vnd.oasis.opendocument.graphics-template"},
	{Key: string(MimeTypeKeyOdf), Name: "odf(OpenDocument Formula)", Kind: "application/vnd.oasis.opendocument.formula"},
	{Key: string(MimeTypeKeyOdc), Name: "odc(OpenDocument Chart)", Kind: "application/vnd.oasis.opendocument.chart"},
	{Key: string(MimeTypeKeySxc), Name: "sxc(Sun XML Calc)", Kind: "application/vnd.sun.xml.calc"},
	{Key: string(MimeTypeKeyPdf), Name: "PDF", Kind: "application/pdf"},
	{Key: string(MimeTypeKeyFdf), Name: "FDF", Kind: "application/vnd.fdf"},
	{Key: string(MimeTypeKeyOleStorage), Name: "OLE Storage", Kind: "application/x-ole-storage"},
	{Key: string(MimeTypeKeyMsi), Name: "MSI", Kind: "application/x-ms-installer"},
	{Key: string(MimeTypeKeyAaf), Name: "AAF", Kind: "application/octet-stream;aaf"},
	{Key: string(MimeTypeKeyMsg), Name: "Outlook", Kind: "application/vnd.ms-outlook"},
	{Key: string(MimeTypeKeyXls), Name: "Excel(xls)", Kind: "application/vnd.ms-excel"},
	{Key: string(MimeTypeKeyPub), Name: "Publisher(pub)", Kind: "application/vnd.ms-publisher"},
	{Key: string(MimeTypeKeyPpt), Name: "Power Point(ppt)", Kind: "application/vnd.ms-powerpoint"},
	{Key: string(MimeTypeKeyDoc), Name: "Word(doc)", Kind: "application/msword"},
	{Key: string(MimeTypeKeyPs), Name: "PostScript", Kind: "application/postscript"},
	{Key: string(MimeTypeKeyPsd), Name: "Photoshop", Kind: "image/vnd.adobe.photoshop"},
	{Key: string(MimeTypeKeyP7s), Name: "PKCS7 Signature", Kind: "application/pkcs7-signature"},
	{Key: string(MimeTypeKeyOgg), Name: "OGG", Kind: "application/ogg"},
	{Key: string(MimeTypeKeyOga), Name: "OGG 音声", Kind: "audio/ogg"},
	{Key: string(MimeTypeKeyOgv), Name: "OGG 動画", Kind: "video/ogg"},
	{Key: string(MimeTypeKeyPng), Name: "PNG 画像", Kind: "image/png"},
	{Key: string(MimeTypeKeyPngApng), Name: "APNG 画像", Kind: "image/vnd.mozilla.apng"},
	{Key: string(MimeTypeKeyJpg), Name: "JPEG 画像", Kind: "image/jpeg"},
	{Key: string(MimeTypeKeyJxl), Name: "JPEG XL 画像", Kind: "image/jxl"},
	{Key: string(MimeTypeKeyJp2), Name: "JPEG 2000 画像", Kind: "image/jp2"},
	{Key: string(MimeTypeKeyJpf), Name: "JPX 画像", Kind: "image/jpx"},
	{Key: string(MimeTypeKeyJpm), Name: "JPM 画像", Kind: "image/jpm"},
	{Key: string(MimeTypeKeyJxs), Name: "JXS 画像", Kind: "image/jxs"},
	{Key: string(MimeTypeKeyGif), Name: "GIF 画像", Kind: "image/gif"},
	{Key: string(MimeTypeKeyWebp), Name: "WebP 画像", Kind: "image/webp"},
	{Key: string(MimeTypeKeyExe), Name: "Portable Executable", Kind: "application/vnd.microsoft.portable-executable"},
	{Key: string(MimeTypeKeyElf), Name: "ELF", Kind: "application/x-elf"},
	{Key: string(MimeTypeKeyObject), Name: "Object", Kind: "application/x-object"},
	{Key: string(MimeTypeKeyExecutable), Name: "実行可能ファイル", Kind: "application/x-executable"},
	{Key: string(MimeTypeKeySharedlib), Name: "Shared Library", Kind: "application/x-sharedlib"},
	{Key: string(MimeTypeKeyCoredump), Name: "Core Dump", Kind: "application/x-coredump"},
	{Key: string(MimeTypeKeyArchive), Name: "Archive", Kind: "application/x-archive"},
	{Key: string(MimeTypeKeyDeb), Name: "Debian Package", Kind: "application/vnd.debian.binary-package"},
	{Key: string(MimeTypeKeyTar), Name: "Tar", Kind: "application/x-tar"},
	{Key: string(MimeTypeKeyXar), Name: "Xar", Kind: "application/x-xar"},
	{Key: string(MimeTypeKeyBz2), Name: "Bzip2", Kind: "application/x-bzip2"},
	{Key: string(MimeTypeKeyFits), Name: "FITS", Kind: "application/fits"},
	{Key: string(MimeTypeKeyTiff), Name: "TIFF 画像", Kind: "image/tiff"},
	{Key: string(MimeTypeKeyBmp), Name: "BMP 画像", Kind: "image/bmp"},
	{Key: string(MimeTypeKeyIcon), Name: "Icon", Kind: "image/x-icon"},
	{Key: string(MimeTypeKeyMpeg), Name: "MPEG 音声", Kind: "audio/mpeg"},
	{Key: string(MimeTypeKeyFlac), Name: "FLAC 音声", Kind: "audio/flac"},
	{Key: string(MimeTypeKeyMidi), Name: "MIDI 音声", Kind: "audio/midi"},
	{Key: string(MimeTypeKeyApe), Name: "APE 音声", Kind: "audio/ape"},
	{Key: string(MimeTypeKeyMpc), Name: "Musepack 音声", Kind: "audio/musepack"},
	{Key: string(MimeTypeKeyAmr), Name: "AMR 音声", Kind: "audio/amr"},
	{Key: string(MimeTypeKeyWav), Name: "WAV 音声", Kind: "audio/wav"},
	{Key: string(MimeTypeKeyAiff), Name: "AIFF 音声", Kind: "audio/aiff"},
	{Key: string(MimeTypeKeyAu), Name: "AU 音声", Kind: "audio/basic"},
	{Key: string(MimeTypeKeyMpegVideo), Name: "MPEG 動画", Kind: "video/mpeg"},
	{Key: string(MimeTypeKeyMov), Name: "QuickTime 動画(mov)", Kind: "video/quicktime"},
	{Key: string(MimeTypeKeyMqv), Name: "QuickTime 動画(mqv)", Kind: "video/quicktime"},
	{Key: string(MimeTypeKeyMp4), Name: "MP4 動画", Kind: "video/mp4"},
	{Key: string(MimeTypeKeyWebm), Name: "WebM 動画", Kind: "video/webm"},
	{Key: string(MimeTypeKey3gp), Name: "3GPP 動画", Kind: "video/3gpp"},
	{Key: string(MimeTypeKey3g2), Name: "3GPP2 動画", Kind: "video/3gpp2"},
	{Key: string(MimeTypeKeyAvi), Name: "AVI 動画", Kind: "video/x-msvideo"},
	{Key: string(MimeTypeKeyFlv), Name: "FLV 動画", Kind: "video/x-flv"},
	{Key: string(MimeTypeKeyMkv), Name: "Matroska 動画", Kind: "video/x-matroska"},
	{Key: string(MimeTypeKeyAsf), Name: "ASF 動画", Kind: "video/x-ms-asf"},
	{Key: string(MimeTypeKeyAac), Name: "AAC 音声", Kind: "audio/aac"},
	{Key: string(MimeTypeKeyVoc), Name: "VOC 音声", Kind: "audio/x-unknown"},
	{Key: string(MimeTypeKeyMp4Audio), Name: "MP4 音声", Kind: "audio/mp4"},
	{Key: string(MimeTypeKeyM4a), Name: "M4A 音声", Kind: "audio/x-m4a"},
	{Key: string(MimeTypeKeyM3u), Name: "M3U", Kind: "application/vnd.apple.mpegurl"},
	{Key: string(MimeTypeKeyM4v), Name: "M4V 動画", Kind: "video/x-m4v"},
	{Key: string(MimeTypeKeyRmvb), Name: "RealMedia VBR", Kind: "application/vnd.rn-realmedia-vbr"},
	{Key: string(MimeTypeKeyGz), Name: "Gzip", Kind: "application/gzip"},
	{Key: string(MimeTypeKeyClass), Name: "Java Applet", Kind: "application/x-java-applet"},
	{Key: string(MimeTypeKeySwf), Name: "Shockwave Flash", Kind: "application/x-shockwave-flash"},
	{Key: string(MimeTypeKeyCrx), Name: "Chrome Extension", Kind: "application/x-chrome-extension"},
	{Key: string(MimeTypeKeyTtf), Name: "TTF", Kind: "font/ttf"},
	{Key: string(MimeTypeKeyWoff), Name: "WOFF", Kind: "font/woff"},
	{Key: string(MimeTypeKeyWoff2), Name: "WOFF2", Kind: "font/woff2"},
	{Key: string(MimeTypeKeyOtf), Name: "OTF", Kind: "font/otf"},
	{Key: string(MimeTypeKeyTtc), Name: "TTC", Kind: "font/collection"},
	{Key: string(MimeTypeKeyEot), Name: "MS Font Object", Kind: "application/vnd.ms-fontobject"},
	{Key: string(MimeTypeKeyWasm), Name: "WebAssembly", Kind: "application/wasm"},
	{Key: string(MimeTypeKeyShx), Name: "SHX", Kind: "application/vnd.shx"},
	{Key: string(MimeTypeKeyShp), Name: "SHP", Kind: "application/vnd.shp"},
	{Key: string(MimeTypeKeyDbf), Name: "DBF", Kind: "application/x-dbf"},
	{Key: string(MimeTypeKeyDcm), Name: "DICOM", Kind: "application/dicom"},
	{Key: string(MimeTypeKeyRar), Name: "RAR", Kind: "application/x-rar-compressed"},
	{Key: string(MimeTypeKeyDjvu), Name: "DjVu", Kind: "image/vnd.djvu"},
	{Key: string(MimeTypeKeyMobi), Name: "Mobipocket", Kind: "application/x-mobipocket-ebook"},
	{Key: string(MimeTypeKeyLit), Name: "MS Reader", Kind: "application/x-ms-reader"},
	{Key: string(MimeTypeKeyBpg), Name: "BPG 画像", Kind: "image/bpg"},
	{Key: string(MimeTypeKeySqlite), Name: "SQLite3", Kind: "application/vnd.sqlite3"},
	{Key: string(MimeTypeKeyDwg), Name: "DWG", Kind: "image/vnd.dwg"},
	{Key: string(MimeTypeKeyNes), Name: "SNES ROM", Kind: "application/vnd.nintendo.snes.rom"},
	{Key: string(MimeTypeKeyLnk), Name: "MS Shortcut", Kind: "application/x-ms-shortcut"},
	{Key: string(MimeTypeKeyMacho), Name: "Mach-O", Kind: "application/x-mach-binary"},
	{Key: string(MimeTypeKeyQcp), Name: "QCELP 音声", Kind: "audio/qcelp"},
	{Key: string(MimeTypeKeyIcns), Name: "ICNS 画像", Kind: "image/x-icns"},
	{Key: string(MimeTypeKeyHeic), Name: "HEIC 画像", Kind: "image/heic"},
	{Key: string(MimeTypeKeyHeicSequence), Name: "HEIC Sequence 画像", Kind: "image/heic-sequence"},
	{Key: string(MimeTypeKeyHeif), Name: "HEIF 画像", Kind: "image/heif"},
	{Key: string(MimeTypeKeyHeifSequence), Name: "HEIF Sequence 画像", Kind: "image/heif-sequence"},
	{Key: string(MimeTypeKeyHdr), Name: "HDR 画像", Kind: "image/vnd.radiance"},
	{Key: string(MimeTypeKeyMrc), Name: "MARC", Kind: "application/marc"},
	{Key: string(MimeTypeKeyMdb), Name: "MS Access(mdb)", Kind: "application/x-msaccess"},
	{Key: string(MimeTypeKeyAccdb), Name: "MS Access(accdb)", Kind: "application/x-msaccess"},
	{Key: string(MimeTypeKeyZst), Name: "Zstandard", Kind: "application/zstd"},
	{Key: string(MimeTypeKeyCab), Name: "MS CAB", Kind: "application/vnd.ms-cab-compressed"},
	{Key: string(MimeTypeKeyRpm), Name: "RPM", Kind: "application/x-rpm"},
	{Key: string(MimeTypeKeyXz), Name: "XZ", Kind: "application/x-xz"},
	{Key: string(MimeTypeKeyLz), Name: "Lzip", Kind: "application/lzip"},
	{Key: string(MimeTypeKeyTorrent), Name: "BitTorrent", Kind: "application/x-bittorrent"},
	{Key: string(MimeTypeKeyCpio), Name: "CPIO", Kind: "application/x-cpio"},
	{Key: string(MimeTypeKeyTzif), Name: "TZIF", Kind: "application/tzif"},
	{Key: string(MimeTypeKeyXcf), Name: "GIMP XCF", Kind: "image/x-xcf"},
	{Key: string(MimeTypeKeyPat), Name: "GIMP PAT", Kind: "image/x-gimp-pat"},
	{Key: string(MimeTypeKeyGbr), Name: "GIMP GBR", Kind: "image/x-gimp-gbr"},
	{Key: string(MimeTypeKeyGlb), Name: "GLTF Binary", Kind: "model/gltf-binary"},
	{Key: string(MimeTypeKeyAvif), Name: "AVIF 画像", Kind: "image/avif"},
	{Key: string(MimeTypeKeyCabInstallshield), Name: "Installshield CAB", Kind: "application/x-installshield"},
	{Key: string(MimeTypeKeyJxr), Name: "JXR 画像", Kind: "image/jxr"},
	{Key: string(MimeTypeKeyTxt), Name: "プレーンテキスト", Kind: "text/plain"},
	{Key: string(MimeTypeKeyHTML), Name: "HTML", Kind: "text/html"},
	{Key: string(MimeTypeKeySvg), Name: "SVG", Kind: "image/svg+xml"},
	{Key: string(MimeTypeKeyXML), Name: "XML", Kind: "text/xml"},
	{Key: string(MimeTypeKeyRss), Name: "RSS", Kind: "application/rss+xml"},
	{Key: string(MimeTypeKeyAtom), Name: "Atom", Kind: "application/atom+xml"},
	{Key: string(MimeTypeKeyX3d), Name: "X3D", Kind: "model/x3d+xml"},
	{Key: string(MimeTypeKeyKml), Name: "Google Earth KML", Kind: "application/vnd.google-earth.kml+xml"},
	{Key: string(MimeTypeKeyXlf), Name: "XLIFF", Kind: "application/x-xliff+xml"},
	{Key: string(MimeTypeKeyDae), Name: "Collada", Kind: "model/vnd.collada+xml"},
	{Key: string(MimeTypeKeyGml), Name: "GML", Kind: "application/gml+xml"},
	{Key: string(MimeTypeKeyGpx), Name: "GPX", Kind: "application/gpx+xml"},
	{Key: string(MimeTypeKeyTcx), Name: "Garmin TCX", Kind: "application/vnd.garmin.tcx+xml"},
	{Key: string(MimeTypeKeyAmf), Name: "AMF", Kind: "application/x-amf"},
	{Key: string(MimeTypeKey3mf), Name: "3D Manufacturing 3D Model", Kind: "application/vnd.ms-package.3dmanufacturing-3dmodel+xml"},
	{Key: string(MimeTypeKeyXfdf), Name: "Adobe XFDF", Kind: "application/vnd.adobe.xfdf"},
	{Key: string(MimeTypeKeyOwl), Name: "OWL", Kind: "application/owl+xml"},
	{Key: string(MimeTypeKeyPhp), Name: "PHP", Kind: "text/x-php"},
	{Key: string(MimeTypeKeyJs), Name: "JavaScript", Kind: "application/javascript"},
	{Key: string(MimeTypeKeyLua), Name: "Lua", Kind: "text/x-lua"},
	{Key: string(MimeTypeKeyPl), Name: "Perl", Kind: "text/x-perl"},
	{Key: string(MimeTypeKeyPy), Name: "Python", Kind: "text/x-python"},
	{Key: string(MimeTypeKeyJSON), Name: "JSON", Kind: "application/json"},
	{Key: string(MimeTypeKeyGeojson), Name: "GeoJSON", Kind: "application/geo+json"},
	{Key: string(MimeTypeKeyHar), Name: "HAR", Kind: "application/json;har"},
	{Key: string(MimeTypeKeyNdjson), Name: "NDJSON", Kind: "application/x-ndjson"},
	{Key: string(MimeTypeKeyRtf), Name: "RTF", Kind: "text/rtf"},
	{Key: string(MimeTypeKeySrt), Name: "SubRip", Kind: "application/x-subrip"},
	{Key: string(MimeTypeKeyTcl), Name: "Tcl", Kind: "text/x-tcl"},
	{Key: string(MimeTypeKeyCsv), Name: "CSV", Kind: "text/csv"},
	{Key: string(MimeTypeKeyTsv), Name: "TSV", Kind: "text/tab-separated-values"},
	{Key: string(MimeTypeKeyVcf), Name: "vCard", Kind: "text/vcard"},
	{Key: string(MimeTypeKeyIcs), Name: "iCalendar", Kind: "text/calendar"},
	{Key: string(MimeTypeKeyWarc), Name: "Warc", Kind: "application/warc"},
	{Key: string(MimeTypeKeyVtt), Name: "VTT", Kind: "text/vtt"},
}

// ManageMimeType マイムタイプ管理サービス。
type ManageMimeType struct {
	DB store.Store
}

// CreateMimeType マイムタイプを作成する。
func (m *ManageMimeType) CreateMimeType(
	ctx context.Context,
	name, key, kind string,
) (entity.MimeType, error) {
	p := parameter.CreateMimeTypeParam{
		Name: name,
		Key:  key,
		Kind: kind,
	}
	e, err := m.DB.CreateMimeType(ctx, p)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to create mime type: %w", err)
	}
	return e, nil
}

// CreateMimeTypes マイムタイプを複数作成する。
func (m *ManageMimeType) CreateMimeTypes(
	ctx context.Context, ps []parameter.CreateMimeTypeParam,
) (int64, error) {
	es, err := m.DB.CreateMimeTypes(ctx, ps)
	if err != nil {
		return 0, fmt.Errorf("failed to create mime types: %w", err)
	}
	return es, nil
}

// UpdateMimeType マイムタイプを更新する。
func (m *ManageMimeType) UpdateMimeType(
	ctx context.Context, id uuid.UUID, name, key, kind string,
) (entity.MimeType, error) {
	p := parameter.UpdateMimeTypeParams{
		Name: name,
		Key:  key,
		Kind: kind,
	}
	e, err := m.DB.UpdateMimeType(ctx, id, p)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to update mime type: %w", err)
	}
	return e, nil
}

// DeleteMimeType マイムタイプを削除する。
func (m *ManageMimeType) DeleteMimeType(ctx context.Context, id uuid.UUID) (int64, error) {
	c, err := m.DB.DeleteMimeType(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete mime type: %w", err)
	}
	return c, nil
}

// PluralDeleteMimeTypes マイムタイプを複数削除する。
func (m *ManageMimeType) PluralDeleteMimeTypes(ctx context.Context, ids []uuid.UUID) (int64, error) {
	c, err := m.DB.PluralDeleteMimeTypes(ctx, ids)
	if err != nil {
		return 0, fmt.Errorf("failed to plural delete mime types: %w", err)
	}
	return c, nil
}

// FindMimeTypeByID マイムタイプをIDで取得する。
func (m *ManageMimeType) FindMimeTypeByID(
	ctx context.Context,
	id uuid.UUID,
) (entity.MimeType, error) {
	e, err := m.DB.FindMimeTypeByID(ctx, id)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to find mime type by id: %w", err)
	}
	return e, nil
}

// FindMimeTypeByKey マイムタイプをキーで取得する。
func (m *ManageMimeType) FindMimeTypeByKey(ctx context.Context, key string) (entity.MimeType, error) {
	e, err := m.DB.FindMimeTypeByKey(ctx, key)
	if err != nil {
		return entity.MimeType{}, fmt.Errorf("failed to find mime type by key: %w", err)
	}
	return e, nil
}

// GetMimeTypes マイムタイプを取得する。
func (m *ManageMimeType) GetMimeTypes(
	ctx context.Context,
	whereSearchName string,
	order parameter.MimeTypeOrderMethod,
	pg parameter.Pagination,
	limit parameter.Limit,
	cursor parameter.Cursor,
	offset parameter.Offset,
	withCount parameter.WithCount,
) (store.ListResult[entity.MimeType], error) {
	wc := store.WithCountParam{
		Valid: bool(withCount),
	}
	var np store.NumberedPaginationParam
	var cp store.CursorPaginationParam
	where := parameter.WhereMimeTypeParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	switch pg {
	case parameter.NumberedPagination:
		np = store.NumberedPaginationParam{
			Valid:  true,
			Offset: entity.Int{Int64: int64(offset), Valid: true},
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.CursorPagination:
		cp = store.CursorPaginationParam{
			Valid:  true,
			Cursor: string(cursor),
			Limit:  entity.Int{Int64: int64(limit), Valid: true},
		}
	case parameter.NonePagination:
	}
	r, err := m.DB.GetMimeTypes(ctx, where, order, np, cp, wc)
	if err != nil {
		return store.ListResult[entity.MimeType]{}, fmt.Errorf("failed to get mime types: %w", err)
	}
	return r, nil
}

// GetMimeTypesCount マイムタイプの数を取得する。
func (m *ManageMimeType) GetMimeTypesCount(
	ctx context.Context,
	whereSearchName string,
) (int64, error) {
	p := parameter.WhereMimeTypeParam{
		WhereLikeName: whereSearchName != "",
		SearchName:    whereSearchName,
	}
	c, err := m.DB.CountMimeTypes(ctx, p)
	if err != nil {
		return 0, fmt.Errorf("failed to get mime types count: %w", err)
	}
	return c, nil
}

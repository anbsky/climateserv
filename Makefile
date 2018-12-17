.:
	go build

clean:
	go clean

install_service:
	sudo cp sds_reader.service /lib/systemd/system/ \
		&& sudo chmod 755 /lib/systemd/system/sds_reader.service \
		&& sudo systemctl enable sds_reader.service \
		&& sudo systemctl start sds_reader

disable_service:
	sudo systemctl disable sds_reader

.PHONY: run
run:
	./sds_reader -i 10 -d /dev/ttyUSB0 -p 2510 -i 10

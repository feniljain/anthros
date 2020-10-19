rm *.pem

#1. Generate CA's private key and self-signed certificate
openssl req -x509 -newkey rsa:4096 -days 365 -keyout ca-key.pem -out ca-cert.pem -subj "/C=IN/ST=Gujarat/L=Surat/O=DSC/OU=DSC-VIT/ CN=*.dscvit.com/emailAddress=fkjainco@gmail.com"

echo "CA's self-signed certificate"
openssl x509 -in ca-cert.pem -noout -text

#2. Generate web server's private key and certificate signing request (CSR)

#3. Use CA's private key to sign web server's CSR and get back the signed certificate

# MKV Production Deployment

## Quick Start

```bash
# 1. Deploy the system
./deploy.sh

# 2. Change default password (IMPORTANT!)
./bin/mkv_client.sh change "fabric_production_password" "your_secure_password"

# 3. Test the system
./bin/mkv_client.sh test "your_secure_password"
```

## Directory Structure

```
deployment/
├── bin/                    # Binaries
│   ├── peer               # Fabric peer with MKV
│   ├── mkv-api-server     # MKV API server
│   └── mkv_client.sh      # MKV client tools
├── lib/                   # Libraries
│   └── libmkv.so         # MKV encryption library
├── config/               # Configuration files
├── scripts/              # Helper scripts
├── mkv-keys/            # MKV key files (persistent)
├── logs/                # Application logs
├── .env                 # Environment variables
├── docker-compose-mkv.yml # Docker compose file
└── deploy.sh            # Quick deployment script
```

## Security Notes

1. **Change default passwords immediately**
2. **Secure the MKV API endpoint**
3. **Regular password rotation**
4. **Monitor logs for security events**
5. **Backup key files regularly**

## Monitoring

```bash
# Health check
curl http://localhost:9876/api/v1/health

# System status
./bin/mkv_client.sh status

# View logs
docker-compose -f docker-compose-mkv.yml logs -f
```

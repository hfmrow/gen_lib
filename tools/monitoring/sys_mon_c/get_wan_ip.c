// get_wan_ip.c

#include "get_wan_ip.h"
#include "file_func.h"

// get numeric ipv4 from adress, i.e "google.com".
char *nslookup(char *hostname, char *hostaddr)
{
	struct addrinfo* res;
	struct sockaddr_in* saddr;

	if (0 != getaddrinfo(hostname, NULL, NULL, &res)) {
		char error[256] = {0};
		sprintf(error, "Unable to resolve: %s\n", hostname);
		internal_error_set(error);
		goto end;
	}

	saddr = (struct sockaddr_in*)res->ai_addr;
	hostaddr = inet_ntoa(saddr->sin_addr);

end:
	return hostaddr;
}

// get public adress ipv4 (wan), this one use a single STUN servers
// "stun_server:port"
bool get_wan_ip(char *stun_server_colon_port, short local_port, char *data)
{
	internal_error_clear();

	struct sockaddr_in servaddr;
	struct sockaddr_in localaddr;
	unsigned char buf[MAXLINE];
	int sockfd, i;
	unsigned char bindingReq[20] = {0};
	char *stun_server_ip;
	char sep[] = {":"};
	short attr_type;
	short attr_length;
	short port;
	short n;
	bool ret = false;

	char *tmp_str = strdup(stun_server_colon_port);
	stun_server_ip = strtok(tmp_str, sep);
	short stun_server_port = atoi(strtok(NULL, sep));

	// Get adress
	stun_server_ip = nslookup(stun_server_ip, data);
	if (ERROR_IS_SET)
		goto end;

	// create socket
	sockfd = socket(AF_INET, SOCK_DGRAM, 0); // UDP
	// configure socket timout for recvfrom
	struct timeval timeout;
	timeout.tv_sec = 5;
	timeout.tv_usec = 0;
	setsockopt(sockfd, SOL_SOCKET, SO_RCVTIMEO, &timeout, sizeof(timeout));

	// server
	bzero(&servaddr, sizeof(servaddr));
	servaddr.sin_family = AF_INET;
	inet_pton(AF_INET, stun_server_ip, &servaddr.sin_addr);
	servaddr.sin_port = htons(stun_server_port);

	// local
	bzero(&localaddr, sizeof(localaddr));
	localaddr.sin_family = AF_INET;
	localaddr.sin_port = htons(local_port);

	n = bind(sockfd,(struct sockaddr *)&localaddr,sizeof(localaddr));

	// first bind
	*(short *)(&bindingReq[0]) = htons(0x0001);		// stun_method
	*(short *)(&bindingReq[2]) = htons(0x0000);		// msg_length
	*(int *)(&bindingReq[4])   = htonl(0x2112A442);	// magic cookie

	*(int *)(&bindingReq[8]) = htonl(0x63c7117e);	// transacation ID
	*(int *)(&bindingReq[12])= htonl(0x0714278f);
	*(int *)(&bindingReq[16])= htonl(0x5ded3221);

	n = sendto(
	        sockfd,
	        bindingReq,
	        sizeof(bindingReq),
	        0,
	        (struct sockaddr *)&servaddr,
	        sizeof(servaddr));
	if (n == -1) {
		internal_error_set("sendto error");
		goto end;
	}
	// time wait
	sleep(1);
	n = recvfrom(sockfd, buf, MAXLINE, 0, NULL,0);
	if (n < 0) {
		if (errno == EWOULDBLOCK) {
			internal_error_set("recvfrom timeout");
			goto end;
		} else {
			internal_error_set("recvfrom socket error");
			goto end;
		}
	}

	if (*(short *)(&buf[0]) == htons(0x0101)) {
		// parse XOR
		n = htons(*(short *)(&buf[2]));
		i = 20;
		while(i<sizeof(buf)) {
			attr_type = htons(*(short *)(&buf[i]));
			attr_length = htons(*(short *)(&buf[i+2]));
			if (attr_type == 0x0020) {
				// parse : port, IP
				port = ntohs(*(short *)(&buf[i+6]));
				port ^= 0x2112;
				sprintf(data,
				        "%d.%d.%d.%d",
				        buf[i+8]^0x21,
				        buf[i+9]^0x12,
				        buf[i+10]^0xA4,
				        buf[i+11]^0x42);
				break;
			}
			i += (4  + attr_length);
		}
	}

	ret = true;
end:
	free(tmp_str);
	close(sockfd);
	return ret;
}
/*
 * some misses to acquire real wan adress using stun srvers has occured,
 * i will try a new approach using 'c' curl library. I don't want to use
 * golang' embedded http lib because it take lot of memory that is useless
 * for a simple http get request.
 */
// get public adress ipv4 (wan), this one use a list of STUN servers
// "stun_server:port"
bool get_wan_adress_from_list(char *list[], int rows_count, char *data)
{
	char error[256];
	for (int i = 0; i < rows_count; i++) {
		printf("adress: %s\n", list[i]);
		if (get_wan_ip(list[i], 8888, data)) {
			return true;
		}
		printf("%s\n", internal_error_get(error));
	}
	return false;
}

// return public adress ipv4 (wan), this one use a list of STUN servers
// "stun_server:port"
char *return_wan_ip(char *stun_server_colon_port, short local_port)
{
	char data[64];

	if (get_wan_ip(stun_server_colon_port, 8888, data))
		return strdup(data);
	else
		return NULL;
}

/*
 * Using curl to perform HTTP GET
 * UBUNTU dev lib: libcurl4-gnutls-dev
 * LINKER option: -lcurl
 */
#include <curl/curl.h>

size_t write_memory_callback(void *contents, size_t size, size_t nmemb, void *userp)
{
	size_t realsize = size * nmemb;
	http_get_memory_struct *mem = ( http_get_memory_struct *)userp;

	char *ptr = realloc(mem->memory, mem->size + realsize + 1);
	if(ptr == NULL) {
		internal_error_set("error: not enough memory");
		return 0;
	}

	mem->memory = ptr;
	memcpy(&(mem->memory[mem->size]), contents, realsize);
	mem->size += realsize;
	mem->memory[mem->size] = 0;

	return realsize;
}

char *curl_http_get(char *adress, char *data)
{
	CURL *curl_handle;
	CURLcode res;

	http_get_memory_struct chunk;
	chunk.memory = malloc(1);
	chunk.size = 0;

	curl_handle = curl_easy_init();
	if(curl_handle) {
		curl_easy_setopt(curl_handle, CURLOPT_URL, adress);
		//curl_easy_setopt(curl_handle, CURLOPT_FOLLOWLOCATION, 1L);
		//curl_easy_setopt(curl_handle, CURLOPT_USERAGENT, "libcurl-agent/1.0");
		curl_easy_setopt(curl_handle, CURLOPT_WRITEFUNCTION, write_memory_callback);
		curl_easy_setopt(curl_handle, CURLOPT_WRITEDATA, (void *)&chunk);

		res = curl_easy_perform(curl_handle);
		if(res != CURLE_OK) {
			sprintf(ERROR_MESSAGE, "%s\n", curl_easy_strerror(res));
			internal_error_set(ERROR_MESSAGE);

			curl_easy_cleanup(curl_handle);
			free(chunk.memory);
			return NULL;
		} else {
			//printf("Size: %lu\n", (unsigned long)chunk.size);
			//printf("Data: %s\n", chunk.memory);
			memcpy(data, chunk.memory, chunk.size);
		}
		curl_easy_cleanup(curl_handle);
		free(chunk.memory);
	}
	return remove_lf(data);
}

// return public adress ipv4 (wan), this one use curl
// 'http get' method. Returned data must be freed.
char *return_wan_ip_http_get(char *adress)
{
	char data[64];
	if (curl_http_get(adress, &data[0]))
		return strdup(data);
	else
		return NULL;
}

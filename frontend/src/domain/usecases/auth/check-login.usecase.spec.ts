import { describe, it, expect, vi } from 'vitest';
import { getJwtPayload } from './get-jwt-payload';
import { verifyTokenValidation } from './check-login.usecase';
import { JwtAuthPayload } from '../../../data/repositories/i-authentication.repo';

// Mock da função getJwtPayload para retornar diferentes cenários
vi.mock('./get-jwt-payload', () => ({
  getJwtPayload: vi.fn(),
}));

describe('verifyTokenValidation', () => {

  it('deve retornar { success: false, message: "Token inválido." } quando o payload do token não contém o campo exp', async () => {
    // Mock para simular um token sem o campo exp
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o campo exp não é um número', async () => {
    // Mock para simular um token com o campo exp não sendo número
    const payload: any = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: 'invalid'
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token expira exatamente agora', async () => {
    // Mock para simular um token que expira exatamente agora
    const currentTime = Math.floor(Date.now() / 1000);
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: currentTime - 1800,
      exp: currentTime
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token contém campos indefinidos', async () => {
    const payload: any = {
      tokenId: undefined,
      isProdutor: undefined,
      session: undefined,
      aud: undefined,
      iat: undefined,
      exp: undefined
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });
  
    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });
  
  it('deve retornar { success: false, message: "Token inválido." } quando o token expira em um futuro muito distante', async () => {
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: 2147483647 // Valor muito alto para exp
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });
  
    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token expira em uma zona de tempo diferente', async () => {
    const currentTime = Math.floor(Date.now() / 1000);
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: currentTime - 1800,
      exp: currentTime + 3600 // Expira em uma hora
    };
    // Mock para alterar a hora do sistema
    vi.spyOn(Date, 'now').mockReturnValue(new Date('2024-08-27T12:00:00Z').getTime());
  
    getJwtPayload.mockResolvedValue({ success: true, data: payload });
  
    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o campo exp é maior que o valor máximo permitido', async () => {
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: 0,
      exp: Number.MAX_SAFE_INTEGER // Valor muito alto para exp
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });
  
    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });
  
  it('deve retornar { success: false, message: "Token inválido." } quando o token contém campos vazios', async () => {
    const payload: JwtAuthPayload = {
      tokenId: '',
      isProdutor: false,
      session: '',
      aud: '',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });
  
    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });
  
  it('deve retornar { success: false, message: "Token inválido." } quando o campo iat é definido para um futuro', async () => {
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) + 3600, // iat no futuro
      exp: Math.floor(Date.now() / 1000) + 7200
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });
  
    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });
  

  it('deve retornar { success: false, message: "Token inválido." } para um token que expira exatamente em uma hora e canTrust é falso', async () => {
    // Mock para simular um token que expira em exatamente uma hora e canTrust falso
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    vi.spyOn(Math, 'random').mockReturnValue(0.01);

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o payload do token contém valores inesperados', async () => {
    // Mock para simular um token com valores inesperados no payload
    const payload: any = {
      tokenId: null,
      isProdutor: null,
      session: null,
      aud: null,
      iat: null,
      exp: null
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o payload do token é vazio', async () => {
    // Mock para simular um payload vazio
    getJwtPayload.mockResolvedValue({ success: true, data: {} });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o payload do token não contém campo iat', async () => {
    // Mock para simular um token sem o campo iat
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando getJwtPayload lança uma exceção', async () => {
    // Mock para simular uma exceção lançada pela função getJwtPayload
    getJwtPayload.mockRejectedValue(new Error('Erro inesperado'));

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o campo exp é menor que o iat', async () => {
    // Mock para simular um token com exp menor que iat
    const currentTime = Math.floor(Date.now() / 1000);
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: currentTime + 3600,
      exp: currentTime
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } para um token com valores extremos para iat e exp', async () => {
    // Mock para simular um token com valores extremos para iat e exp
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: 0,
      exp: 2147483647 // Valor muito alto para exp
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    // Mock de Math.random para retornar um valor que fará canTrust ser falso
    vi.spyOn(Math, 'random').mockReturnValue(0.01);

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o payload do token contém iat maior que exp', async () => {
    // Mock para simular um token com iat maior que exp
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) + 3600,
      exp: Math.floor(Date.now() / 1000) + 1800
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token é válido mas a hora atual está errada', async () => {
    // Mock para simular uma situação onde a hora atual está incorreta
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    // Mock de Date.now() para retornar um valor que faz o token parecer inválido
    vi.spyOn(Date, 'now').mockReturnValue(Date.now() + 7200 * 1000); // Adiciona 2 horas

    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o payload contém valores negativos para iat e exp', async () => {
    // Mock para simular um token com valores negativos para iat e exp
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: -3600,
      exp: -1800
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token contém dados inválidos no payload', async () => {
    // Mock para simular um token com dados inválidos
    const payload: any = {
      tokenId: '123',
      isProdutor: 'invalid',
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token é válido mas o random retorna 0', async () => {
    // Mock para simular um token válido e random retornando 0
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    vi.spyOn(Math, 'random').mockReturnValue(0);

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: true } quando o token é válido e canTrust é 1', async () => {
    // Mock para simular um token válido e canTrust verdadeiro
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    vi.spyOn(Math, 'random').mockReturnValue(1);

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      success: true
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o payload é um objeto vazio', async () => {
    // Mock para simular um payload vazio
    getJwtPayload.mockResolvedValue({ success: true, data: {} });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token é válido mas não contém campo isProdutor', async () => {
    // Mock para simular um token sem o campo isProdutor
    const payload: JwtAuthPayload = {
      tokenId: '123',
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token é válido mas não contém campo aud', async () => {
    // Mock para simular um token sem o campo aud
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token é válido mas não contém campo tokenId', async () => {
    // Mock para simular um token sem o campo tokenId
    const payload: JwtAuthPayload = {
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token é válido mas não contém campo session', async () => {
    // Mock para simular um token sem o campo session
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token é válido mas contém valores em formato incorreto', async () => {
    // Mock para simular um token com valores em formato incorreto
    const payload: any = {
      tokenId: 123,
      isProdutor: 'true',
      session: 123,
      aud: 123,
      iat: 'invalid',
      exp: 'invalid'
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token contém campos vazios ou nulos', async () => {
    // Mock para simular um token com campos vazios ou nulos
    const payload: JwtAuthPayload = {
      tokenId: '',
      isProdutor: null,
      session: '',
      aud: '',
      iat: null,
      exp: null
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: false, message: "Token inválido." } quando o token expira em exatamente uma hora e canTrust é falso', async () => {
    // Mock para simular um token que expira em exatamente uma hora e canTrust falso
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    vi.spyOn(Math, 'random').mockReturnValue(0.01);

    const result = await verifyTokenValidation();
    expect(result).toEqual({
      "message": "Token inválido.",
      "success": false,
    });
  });

  it('deve retornar { success: true } quando o token é válido e canTrust é verdadeiro em múltiplas iterações', async () => {
    // Mock para simular um token válido e canTrust verdadeiro em múltiplas iterações
    const payload: JwtAuthPayload = {
      tokenId: '123',
      isProdutor: false,
      session: 'session123',
      aud: 'aud123',
      iat: Math.floor(Date.now() / 1000) - 1800,
      exp: Math.floor(Date.now() / 1000) + 3600
    };
    getJwtPayload.mockResolvedValue({ success: true, data: payload });

    vi.spyOn(Math, 'random').mockReturnValue(0.95);

    // Testar múltiplas iterações para garantir consistência
    for (let i = 0; i < 10; i++) {
      const result = await verifyTokenValidation();
      expect(result).toEqual({
        success: true
      });
    }
  });

});
